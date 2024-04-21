package agents

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func HandleAgentResponse(messages <-chan amqp.Delivery) {
	if err := InitAgents(); err != nil {
		logrus.Errorf("Не удалось загрузить список агентов: %s", err.Error())
		return
	}

	for message := range messages {
		logrus.Debugf("[rabbitmq] тип полученного сообщения \"%s\" с телом: \"%s\"", message.Type, message.Body)

		if message.Type == "answer" {
			HandleAnswer(message)
		}
		if message.Type == "processed" {
			HandleProcessed(message)
		}
		if message.Type == "error" {
			HandleError(message)
		}
		if message.Type == "ping" {
			HandlePing(message)
		}
	}
}
