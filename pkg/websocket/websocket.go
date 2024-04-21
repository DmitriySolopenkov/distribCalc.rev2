package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type WSData struct {
	Action string      `json:"action"`
	Id     interface{} `json:"id"`
	Data   interface{} `json:"data"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // for cors
	}
	clients []*websocket.Conn
)

func Connect(ctx *gin.Context) error {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	clients = append(clients, conn)
	return nil
}

func Broadcast(data WSData) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for index, conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			clients = append(clients[:index], clients[index+1:]...) // remove connection from slices
			logrus.Infof("Websocket %s has been disconnected", conn.RemoteAddr().String())
			continue
		}
	}

	return nil
}
