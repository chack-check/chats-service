package ws

import (
	"github.com/chack-check/chats-service/protousers"
	"github.com/gorilla/websocket"
)

type Client struct {
	id         string
	user       *protousers.UserResponse
	connection *websocket.Conn
	send       chan []byte
	handlers   []ClientMessageHandler
}

func (c *Client) read() {
	defer func() {
		Pool.unregister <- c
		c.connection.Close()
	}()

	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			Pool.unregister <- c
			c.connection.Close()
			break
		}

		sended := false
		for _, handler := range c.handlers {
			responseMessage, err := handler.HandleMessage(message, c)
			if err == nil {
				sended = true
				c.send <- responseMessage
				break
			} else if len(responseMessage) > 0 {
				sended = true
				c.send <- responseMessage
				break
			}
		}

		if !sended {
			c.send <- NewUndefinedMessageBytes()
		}
	}
}

func (c *Client) write() {
	defer func() {
		Pool.unregister <- c
		c.connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.connection.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func NewClient(id string, connection *websocket.Conn) *Client {
	return &Client{
		id:         id,
		connection: connection,
		handlers: []ClientMessageHandler{
			AuthMessageHandler{},
		},
	}
}
