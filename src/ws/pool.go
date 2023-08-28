package ws

import (
	"slices"
)

type ClientsPool struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (p *ClientsPool) Start() {
	for {
		select {
		case conn := <-p.register:
			p.clients[conn] = true
		case conn := <-p.unregister:
			close(conn.send)
			delete(p.clients, conn)
		}
	}
}

func (p *ClientsPool) Send(event []byte, ids []int32) {
	for client, ok := range p.clients {
		if slices.Contains(ids, client.user.Id) && ok {
			client.send <- event
		}
	}
}

var Pool = ClientsPool{}
