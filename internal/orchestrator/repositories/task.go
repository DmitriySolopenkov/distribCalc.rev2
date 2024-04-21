package repositories

import (
	"context"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/database"
	"time"
)

type Task struct {
}

type TaskModel struct {
	TaskID     int       `json:"task_id"`
	Expression string    `json:"expression"`
	Status     string    `json:"status"`
	Answer     string    `json:"answer"`
	AgentID    string    `json:"agent_id"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

const (
	STATUS_CREATED   = "created"
	STATUS_PROCESSED = "processed"
	STATUS_FAIL      = "fail"
	STATUS_COMPLETED = "completed"
)

// Get all tasks in database
func (t *Task) GetAllTasks(userId int) ([]TaskModel, error) {
	rows, err := database.DB.Query(context.Background(), "SELECT * FROM tasks WHERE user_id = $1 ORDER BY task_id DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []TaskModel{}

	for rows.Next() {
		var task TaskModel
		if err = rows.Scan(&task.TaskID, &task.Expression, &task.Status, &task.Answer, &task.AgentID, &task.UserID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Create new row with task
func (t *Task) Create(expression string, userId int) (int, error) {
	var insertedID int

	query := "INSERT INTO tasks (expression, status, answer, agent_id, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING task_id"
	if err := database.DB.QueryRow(context.Background(), query, expression, STATUS_CREATED, "", "", userId).Scan(&insertedID); err != nil {
		return 0, err
	}

	return insertedID, nil
}

// Find task by id
func (t *Task) GetById(taskId int) (TaskModel, error) {
	var task TaskModel

	query := "SELECT * FROM tasks WHERE task_id = $1"
	if err := database.DB.QueryRow(context.Background(), query, taskId).Scan(&task.TaskID, &task.Expression, &task.Status, &task.Answer, &task.AgentID, &task.UserID, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return task, err
	}

	return task, nil
}

// Update answer for expression by task id
func (t *Task) SetAnswer(taskId int, answer string, status string) error {
	query := "UPDATE tasks SET answer = $1, status = $2, updated_at = $3 WHERE task_id = $4"
	if _, err := database.DB.Exec(context.Background(), query, answer, status, time.Now(), taskId); err != nil {
		return err
	}

	return nil
}

// Update status to processed and set agent_id
func (t *Task) SetProcessed(taskId int, agentId string) error {
	query := "UPDATE tasks SET status = $1, agent_id = $2, updated_at = $3 WHERE task_id = $4"
	if _, err := database.DB.Exec(context.Background(), query, STATUS_PROCESSED, agentId, time.Now(), taskId); err != nil {
		return err
	}

	return nil
}

// Get all tasks for resolve
func (t *Task) GetTasksForResolve() ([]TaskModel, error) {
	sql := "SELECT * FROM tasks WHERE status = $1 AND updated_at < $2"
	timeForResolve := time.Now().Add(-time.Duration(config.Get().AgentResolveTime) * time.Second)
	rows, err := database.DB.Query(context.Background(), sql, STATUS_PROCESSED, timeForResolve)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []TaskModel{}

	for rows.Next() {
		var task TaskModel
		if err = rows.Scan(&task.TaskID, &task.Expression, &task.Status, &task.Answer, &task.AgentID, &task.UserID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update status to created and unset agent_id
func (t *Task) SetCreated(taskId int) error {
	query := "UPDATE tasks SET status = $1, agent_id = $2, updated_at = $3 WHERE task_id = $4"
	if _, err := database.DB.Exec(context.Background(), query, STATUS_CREATED, "", time.Now(), taskId); err != nil {
		return err
	}

	return nil
}

func TaskRepository() *Task {
	return &Task{}
}
