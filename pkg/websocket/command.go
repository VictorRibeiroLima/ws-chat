package websocket

import "fmt"

type Command struct {
	Command string            `json:"command"`
	Data    map[string]string `json:"data"`
}

func (comand *Command) Do(c *Client) {
	switch comand.Command {
	case "SET_NICK":
		{
			nick := comand.Data["nick"]
			message := Message{Nick: "Server", Body: "User '" + c.Nick + "' change nick to '" + nick + "'"}
			payload := Payload{From: c, Message: message}
			c.Nick = nick
			c.Pool.Broadcast <- payload
		}
	case "SEND_MESSAGE":
		{
			message := Message{Nick: c.Nick, Body: comand.Data["message"]}
			payload := Payload{From: c, Message: message}
			c.Pool.Broadcast <- payload
			fmt.Printf("Message Received: %+v\n", message)
		}
	}
}
