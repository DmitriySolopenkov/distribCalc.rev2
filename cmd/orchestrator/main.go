package main

import (
	"flag"
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator"
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/agents"
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/routes"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/database"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/logger"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Expression Calculator
// @version 1.0
// @description Endpoints for expression calculator by Dmitriy Solopenkov
// @BasePath /api/v1
func main() {
	// enable or disable debug mode
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	// initialization a logrus
	logger.Init(debug)
	// parsing .env file to config struct
	config_, err := config.Init()
	if err != nil {
		logrus.Fatalf("failed parse config .env: %s", err.Error())
		return
	}
	logrus.Debug("config was successful loaded")

	// try to connect to database
	if err = database.Init(); err != nil {
		logrus.Fatalf("database connection failed: %s", err.Error())
		return
	}
	if err = orchestrator.PrepareDatabase(); err != nil {
		logrus.Fatalf("database prepare sql failed: %s", err.Error())
		return
	}
	logrus.Debug("database was successful loaded & prepared")

	// try to connect to rabbitmq
	broker, err := rabbitmq.Init(config_.RabbitDSN)
	if err != nil {
		logrus.Fatal("rabbitmq connection failed")
		return
	}
	// try to create a queue for send task to rabbitmq
	if err = broker.InitQueue(config_.RabbitTaskQueue); err != nil {
		logrus.Fatal("rabbitmq fail creation a queue for tasks")
		return
	}
	// try to create a queue for server responses
	if err = broker.InitQueue(config_.RabbitAgentQueue); err != nil {
		logrus.Fatal("rabbitmq fail creation a queue for servers")
		return
	}
	logrus.Debug("rabbitmq successful connected")

	// start listen a responses from agents
	messages, err := broker.ConnQueue(config_.RabbitAgentQueue)
	go agents.HandleAgentResponse(messages)
	go agents.HandleTimeoutAgents()
	go orchestrator.ResolveTasks()

	// initialization a gin
	gin.SetMode(config_.Mode)
	router := gin.Default()
	routes.InitRouter(router)

	logrus.Info("Orchestrator was successful started!")
	// run a server
	if err = router.Run(config_.ServerAddr); err != nil {
		logrus.Fatalf("failed run http server: %s", err.Error())
	}
}
