package repositories

import (
	"context"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"time"

	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/database"
)

type Agent struct {
}

type AgentModel struct {
	AgentID  string    `json:"agent_id"`
	LastPing time.Time `json:"last_ping"`
	Status   string    `json:"status"`
}

const (
	AGENT_CONNECTED    = "connected"
	AGENT_DISCONNECTED = "disconnected"
	AGENT_DELETED      = "deleted"
)

// Get all agents in database
func (a *Agent) GetAllAgents() ([]AgentModel, error) {
	rows, err := database.DB.Query(context.Background(), "SELECT * FROM agents ORDER BY last_ping DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := []AgentModel{}

	for rows.Next() {
		var agent AgentModel
		if err = rows.Scan(&agent.AgentID, &agent.LastPing); err != nil {
			return nil, err
		}

		if time.Now().Add(-time.Duration(config.Get().AgentPing) * time.Second).After(agent.LastPing) {
			agent.Status = AGENT_DISCONNECTED
		} else {
			agent.Status = AGENT_CONNECTED
		}

		agents = append(agents, agent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

// Create new row with agent
func (a *Agent) Create(agentId string) error {
	query := "INSERT INTO agents (agent_id, last_ping) VALUES ($1, $2)"
	if _, err := database.DB.Exec(context.Background(), query, agentId, time.Now()); err != nil {
		return err
	}

	return nil
}

// Update last ping by agent id
func (a *Agent) SetLastPing(agentId string, ping time.Time) error {
	query := "UPDATE agents SET last_ping = $1 WHERE agent_id = $2"
	if _, err := database.DB.Exec(context.Background(), query, ping, agentId); err != nil {
		return err
	}

	return nil
}

// Delete agent by agent id
func (a *Agent) Delete(agentId string) error {
	query := "DELETE FROM agents WHERE agent_id = $1"
	if _, err := database.DB.Exec(context.Background(), query, agentId); err != nil {
		return err
	}

	return nil
}

func AgentRepository() *Agent {
	return &Agent{}
}
