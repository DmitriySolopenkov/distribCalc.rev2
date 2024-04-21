package agent

import (
	"fmt"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/rabbitmq"
	"github.com/Knetic/govaluate"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

func Solver(queueOrchestrator, agentId string, wait int, messages <-chan amqp.Delivery) {
	for message := range messages {
		if message.Type != "task" {
			continue
		}

		logrus.Infof("Received task #%s: %s", message.CorrelationId, message.Body)

		processed := amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(agentId),
			Type:          "processed",
			CorrelationId: message.CorrelationId,
		}
		if err := rabbitmq.Get().SendToQueue(queueOrchestrator, processed); err != nil {
			logrus.Fatalf("Failed sent status processed for task #%s: %s", message.CorrelationId, err.Error())
			continue
		}

		time.Sleep(time.Duration(wait) * time.Second) // wait for response

		expression, err := govaluate.NewEvaluableExpression(string(message.Body))
		if err != nil {
			errorMsg := amqp.Publishing{
				ContentType:   "text/plain",
				Body:          []byte(err.Error()),
				Type:          "error",
				CorrelationId: message.CorrelationId,
			}
			if err := rabbitmq.Get().SendToQueue(queueOrchestrator, errorMsg); err != nil {
				logrus.Fatalf("Failed sent error for task #%s: %s", message.CorrelationId, err.Error())
				continue
			}

			logrus.Errorf("Failed load expression task #%s \"%s\": %s", message.CorrelationId, message.Body, err.Error())
			continue
		}
		result, err := expression.Evaluate(map[string]interface{}{})

		if err != nil {
			errorMsg := amqp.Publishing{
				ContentType:   "text/plain",
				Body:          []byte(err.Error()),
				Type:          "error",
				CorrelationId: message.CorrelationId,
			}
			if err := rabbitmq.Get().SendToQueue(queueOrchestrator, errorMsg); err != nil {
				logrus.Fatalf("Failed sent error for task #%s: %s", message.CorrelationId, err.Error())
				continue
			}

			logrus.Errorf("Failed solve the expression %s for task #%s: %s", message.Body, message.CorrelationId, err.Error())
			continue
		}
		resultByte := []byte(fmt.Sprint(result))

		answer := amqp.Publishing{
			ContentType:   "text/plain",
			Body:          resultByte,
			Type:          "answer",
			CorrelationId: message.CorrelationId,
		}

		if err = rabbitmq.Get().SendToQueue(queueOrchestrator, answer); err != nil {
			logrus.Fatalf("Failed sent answer %s for task #%s: %s", resultByte, message.CorrelationId, err.Error())
			continue
		}

		logrus.Infof("Answer sent for task #%s: %s", message.CorrelationId, resultByte)
	}
}
