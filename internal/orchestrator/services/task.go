package services

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/rabbitmq"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/websocket"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

type Task struct {
	repo *repositories.Task
}

func (t *Task) Create(expression string, userId int) (repositories.TaskModel, error) {
	taskId, err := t.repo.Create(expression, userId)

	if err != nil {
		return repositories.TaskModel{}, err
	}

	// send task to queue of rabbitmq

	message := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          []byte(expression),
		Type:          "task",
		CorrelationId: strconv.Itoa(taskId),
	}

	if err = rabbitmq.Get().SendToQueue(config.Get().RabbitTaskQueue, message); err != nil {
		return repositories.TaskModel{}, err
	}

	// return task model for response

	task := repositories.TaskModel{
		TaskID:     taskId,
		Status:     repositories.STATUS_CREATED,
		Expression: expression,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// send message to websocket

	wsData := websocket.WSData{
		Action: "new_task",
		Id:     taskId,
		Data:   task,
	}
	if err = websocket.Broadcast(wsData); err != nil {
		return repositories.TaskModel{}, err
	}

	return task, nil
}

func (t *Task) SetAnswer(taskId int, answer, status string) error {
	if err := repositories.TaskRepository().SetAnswer(taskId, answer, status); err != nil {
		logrus.Errorf("Failed update a row with task %d: %s", taskId, err.Error())
		return err
	}

	task, err := repositories.TaskRepository().GetById(taskId)
	if err != nil {
		logrus.Errorf("Failed fetch a row with task by taskId #%d: %s", taskId, err.Error())
		return err
	}

	wsData := websocket.WSData{
		Action: "update_task",
		Id:     taskId,
		Data:   task,
	}
	if err = websocket.Broadcast(wsData); err != nil {
		return err
	}

	return nil
}

func (t *Task) SetProcessed(taskId int, agentId string) error {
	if err := repositories.TaskRepository().SetProcessed(taskId, agentId); err != nil {
		logrus.Errorf("Failed update a row with task %d: %s", taskId, err.Error())
		return err
	}

	task, err := repositories.TaskRepository().GetById(taskId)
	if err != nil {
		logrus.Errorf("Failed fetch a row with task by taskId #%d: %s", taskId, err.Error())
		return err
	}

	wsData := websocket.WSData{
		Action: "update_task",
		Id:     taskId,
		Data:   task,
	}
	if err = websocket.Broadcast(wsData); err != nil {
		return err
	}

	return nil
}

func (t *Task) ResolveTask(task repositories.TaskModel) error {
	message := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          []byte(task.Expression),
		Type:          "task",
		CorrelationId: strconv.Itoa(task.TaskID),
	}

	if err := rabbitmq.Get().SendToQueue(config.Get().RabbitTaskQueue, message); err != nil {
		logrus.Fatalf("Failed send message to rabbitmq: %s", err.Error())
		return err
	}

	if err := repositories.TaskRepository().SetCreated(task.TaskID); err != nil {
		logrus.Errorf("Failed update status to created when resolved for task #%d: %s", task.TaskID, err.Error())
		return err
	}

	wsData := websocket.WSData{
		Action: "update_task",
		Id:     task.TaskID,
		Data:   task,
	}
	if err := websocket.Broadcast(wsData); err != nil {
		return err
	}

	return nil
}

// create new task service
func TaskService() *Task {
	return &Task{}
}
