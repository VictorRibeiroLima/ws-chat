package websocket

import "fmt"

var id = 1
var roomSerialId *int = &id

type Room struct {
	ID         int
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Broadcast  chan Payload
}

func NewRoom() *Room {
	id := *roomSerialId
	(*roomSerialId)++
	return &Room{
		ID:         id,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan Payload),
	}
}

func (room *Room) Start() {
	room.listen()
}

func (room *Room) listen() {
	for {
		select {
		case client := <-room.Register:
			{
				room.Clients[client.ID] = client
				for id, client := range room.Clients {
					fmt.Printf("Client with id %s entered room %d \n", id, room.ID)
					client.Conn.WriteJSON(Message{Nick: "Server", Body: "New User Joined:" + client.Nick})
				}
				fmt.Println("Size of Connection Room: ", len(room.Clients))
			}
		case client := <-room.Unregister:
			{
				delete(room.Clients, client.ID)
				fmt.Println("Size of Connection Room: ", len(room.Clients))
				for _, client := range room.Clients {
					fmt.Printf("Client with id %d leaved room %d \n", id, room.ID)
					client.Conn.WriteJSON(Message{Nick: "Server", Body: "User Disconnected:" + client.Nick})
				}
			}
		case payload := <-room.Broadcast:
			{
				for _, client := range room.Clients {
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
}
