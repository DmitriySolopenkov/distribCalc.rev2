package orchestrator

import (
	"context"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/database"
)

func PrepareDatabase() error {

	var sql = []string{
		"set time zone 'Europe/Moscow'",
		"create table if not exists tasks (task_id serial primary key, expression text not null, status varchar(10) not null, answer text not null, agent_id varchar(255), user_id integer, created_at timestamp with time zone default CURRENT_TIMESTAMP, updated_at timestamp with time zone default CURRENT_TIMESTAMP);",
		"create table if not exists agents (agent_id varchar(255) primary key, last_ping timestamp with time zone default CURRENT_TIMESTAMP);",
		"create table if not exists users (user_id serial primary key, login varchar(255), password varchar(255), created_at timestamp with time zone default CURRENT_TIMESTAMP, updated_at timestamp with time zone default CURRENT_TIMESTAMP);",
	}

	for _, query := range sql {
		if _, err := database.DB.Query(context.Background(), query); err != nil {
			return err
		}
	}

	return nil
}
