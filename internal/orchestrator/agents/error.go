package agents

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/services"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
)

func HandleError(message amqp.Delivery) {
	// handle error from agent
	taskId, err := strconv.Atoi(message.CorrelationId)
	if err != nil {
		logrus.Errorf("Получен неверный идентификатор задачи из-за ошибки: %s", message.CorrelationId)
		return
	}

	// set error into database
	if err = services.TaskService().SetAnswer(taskId, string(message.Body), repositories.STATUS_FAIL); err != nil {
		logrus.Errorf("Не удалось обновить строку с задачей %d: %s", taskId, err.Error())
		return
	}

	logrus.Infof("Получена ошибка для %s: %s", message.CorrelationId, message.Body)
}
