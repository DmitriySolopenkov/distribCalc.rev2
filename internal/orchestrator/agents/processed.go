package agents

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/services"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
)

func HandleProcessed(message amqp.Delivery) {
	// обрабатывать статус обработки от агента
	taskId, err := strconv.Atoi(message.CorrelationId)
	if err != nil {
		logrus.Errorf("Получен неверный идентификатор задачи для обработки %s", message.CorrelationId)
		return
	}

	// установить статус обработки в базу данных
	if err = services.TaskService().SetProcessed(taskId, string(message.Body)); err != nil {
		logrus.Errorf("Не удалось обновить строку с задачей %d: %s", taskId, err.Error())
		return
	}

	logrus.Infof("Установлен статус обработан для #%s с помощью агента #%s", message.CorrelationId, message.Body)
}
