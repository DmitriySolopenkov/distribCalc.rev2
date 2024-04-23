package agents

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/websocket"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

var agents = make(map[string]*repositories.AgentModel)

func HandlePing(message amqp.Delivery) {
	// обрабатывать пинг от агента (создать строку или обновить Last_ping)
	timeout := time.Duration(config.Get().AgentTimeout) * time.Second
	if time.Since(message.Timestamp) >= timeout {
		// устаревшее сообщение ping
		logrus.Debugf("Устаревшее сообщение для пинга %s", message.Body)
		return
	}

	agent := string(message.Body)
	_, ok := agents[agent]

	if !ok {
		// это новый агент, создайте строку
		if err := repositories.AgentRepository().Create(agent); err != nil {
			logrus.Fatalf("Не удалось создать нового агента %s: %s", agent, err.Error())
			return
		}

		agents[agent] = &repositories.AgentModel{
			AgentID:  agent,
			LastPing: message.Timestamp,
			Status:   repositories.AGENT_CONNECTED,
		}
		sendToWebsocket(*agents[agent])

		logrus.Infof("Подключен новый агент #%s", agent)
	}

	if err := repositories.AgentRepository().SetLastPing(agent, message.Timestamp); err != nil {
		// ошибка при обновлении базы данных
		logrus.Fatalf("Не удалось обновить последний пинг для %s: %s", agent, err.Error())
		return
	}

	if ok {
		agents[agent].LastPing = message.Timestamp
	}

	logrus.Debugf("Обновить последний пинг для %s", agent)
}

func HandleTimeoutAgents() {
	// каждую секунду проверяйте, какие агенты отключены
	ticker := time.NewTicker(time.Second)
	timeout := time.Duration(config.Get().AgentTimeout) * time.Second
	pingTime := time.Duration(config.Get().AgentPing) * time.Second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for agentID, agent := range agents {
				if time.Now().Add(-timeout).After(agent.LastPing) {
					// агент отключен более 10 минут - удалить из базы
					agents[agentID].Status = repositories.AGENT_DELETED
					sendToWebsocket(*agents[agentID])

					delete(agents, agentID)
					if err := repositories.AgentRepository().Delete(agentID); err != nil {
						logrus.Errorf("Не удалось удалить агент #%s: %s", agent, err.Error())
						continue
					}

					logrus.Infof("Тайм-аут агента #%s (удаление)", agentID)
				}

				if time.Now().Add(-pingTime).After(agent.LastPing) {
					// агент отключен менее 10 минут - установить статус отключен

					if agent.Status != repositories.AGENT_CONNECTED {
						// уже отправленное сообщение о том, что агент отключен
						continue
					}

					agents[agentID].Status = repositories.AGENT_DISCONNECTED
					sendToWebsocket(*agents[agentID])
					logrus.Infof("Агент #%s отключен", agentID)
				} else if agent.Status != repositories.AGENT_CONNECTED {
					// agent has been reconnected

					agents[agentID].Status = repositories.AGENT_CONNECTED
					sendToWebsocket(*agents[agentID])
					logrus.Infof("Агент #%s будет переподключен", agentID)
				}
			}
		}
	}
}

func InitAgents() error {
	// загрузка текущих агентов из базы данных
	agentsDb, err := repositories.AgentRepository().GetAllAgents()
	if err != nil {
		return err
	}

	for _, agent := range agentsDb {
		agents[agent.AgentID] = &agent
	}

	return nil
}

func sendToWebsocket(agent repositories.AgentModel) {
	wsData := websocket.WSData{
		Action: "update_agent",
		Id:     agent.AgentID,
		Data:   agent,
	}
	if err := websocket.Broadcast(wsData); err != nil {
		logrus.Errorf("Не удалось отправить сообщение об агенте в веб-сокет #%s: %s", agent.AgentID, err.Error())
		return
	}
}
