package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/VictorRibeiroLima/ws-server/pkg/util"
	"github.com/VictorRibeiroLima/ws-server/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	nick := r.URL.Query().Get("nick")
	if nick == "" {
		nick = "anonymous"
	}
	param := strings.TrimPrefix(r.URL.Path, "/ws/")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s is not a number", param)
		return
	}
	room := pool.Rooms[id]
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Rom with id '%d' not found", id)
		return
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to upgrade connection")
		return
	}

	client := &websocket.Client{
		ID:   util.RandomString(),
		Nick: nick,
		Conn: conn,
		Room: room,
	}

	room.Register <- client
	client.Read()
}

func listRooms(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	type Room struct {
		ID     int `json:"id"`
		Lenght int `json:"lenght"`
	}
	var rooms []Room
	for _, r := range pool.Rooms {
		room := Room{
			ID:     r.ID,
			Lenght: len(r.Clients),
		}
		rooms = append(rooms, room)
	}
	sort.SliceStable(rooms, func(i, j int) bool {
		return rooms[i].ID < rooms[j].ID
	})
	responseBody, err := json.Marshal(rooms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to recover rooms")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func addRoom(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	room := websocket.NewRoom()
	pool.Register <- room
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Room added")
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		addRoom(pool, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		listRooms(pool, w, r)
	})

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Server started on port 3000")
	setupRoutes()
	http.ListenAndServe(":3000", nil)

}
