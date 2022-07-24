package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Payload
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Payload),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Nick: "Server", Body: "New User Joined:" + client.Nick})
			}
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Nick: "Server", Body: "User Disconnected:" + client.Nick})
			}
			break
		case payload := <-pool.Broadcast:
			for client := range pool.Clients {
				if payload.From != nil && client == payload.From {
					continue
				}
				if err := client.Conn.WriteJSON(payload.Message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
