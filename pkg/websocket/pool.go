package websocket

import "fmt"

type Pool struct {
	Register   chan *Room
	Unregister chan *Room
	Rooms      map[int]*Room
}

func NewPool() *Pool {
	pool := &Pool{
		Register:   make(chan *Room),
		Unregister: make(chan *Room),
		Rooms:      make(map[int]*Room),
	}
	for i := 0; i < 20; i++ {
		room := NewRoom()
		pool.Rooms[room.ID] = room
		go room.Start()
	}
	return pool
}

func (pool *Pool) Start() {
	for {
		select {
		case room := <-pool.Register:
			{
				pool.Rooms[room.ID] = room
				fmt.Printf("New room added with id '%d' \n", room.ID)
			}
		case room := <-pool.Unregister:
			{
				delete(pool.Rooms, room.ID)
				fmt.Printf("Room was '%d' delete \n", room.ID)
			}
		}
	}
}
