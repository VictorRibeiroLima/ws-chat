package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Nick string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Nick string `json:"nick"`
	Body string `json:"body"`
}

type Payload struct {
	From    *Client
	Message Message
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var command Command
		json.Unmarshal(p, &command)
		command.Do(c)
	}
}
