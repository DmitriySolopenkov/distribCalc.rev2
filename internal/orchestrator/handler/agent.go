package handler

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/response"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/websocket"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Agent struct {
	Route *gin.RouterGroup
}

// @Summary Get all agents
// @Tags Agent
// @ID agent-index
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]repositories.AgentModel}
// @Router /agent [get]
func (a *Agent) Index(ctx *gin.Context) {
	agents, err := repositories.AgentRepository().GetAllAgents()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Data(ctx, agents)
}

func (a *Agent) WebSocket(ctx *gin.Context) {
	err := websocket.Connect(ctx)
	if err != nil {
		logrus.Errorln("Wrong websocket connection: %s", err.Error())
		response.BadRequest(ctx, "incorrect websocket")
		return
	}
}
