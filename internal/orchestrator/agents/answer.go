package agents

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/services"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
)

func HandleAnswer(message amqp.Delivery) {
	// handle answer from agent
	taskId, err := strconv.Atoi(message.CorrelationId)
	if err != nil {
		logrus.Errorf("Получен неправильный идентификатор задачи для ответа: %s", message.CorrelationId)
		return
	}

	// set answer into database
	if err = services.TaskService().SetAnswer(taskId, string(message.Body), repositories.STATUS_COMPLETED); err != nil {
		logrus.Errorf("Не удалось обновить строку с задачей %d: %s", taskId, err.Error())
		return
	}

	logrus.Infof("Получен ответ на %s: %s", message.CorrelationId, message.Body)
}
