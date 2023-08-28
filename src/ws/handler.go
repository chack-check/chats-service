package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	client := &Client{
		id:         uuid.Must(uuid.NewV4(), nil).String(),
		connection: conn,
		send:       make(chan []byte),
	}

	Pool.register <- client

	go client.read()
	go client.write()
}
