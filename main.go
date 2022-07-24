package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VictorRibeiroLima/ws-server/pkg/util"
	"github.com/VictorRibeiroLima/ws-server/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}

	client := &websocket.Client{
		ID:   util.RandomString(),
		Nick: "anonymous",
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Server started on port 8080")
	setupRoutes()
	http.ListenAndServe(":8080", nil)

}
