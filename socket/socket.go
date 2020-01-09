package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients []*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	clients = append(clients, ws)
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte(`{"hello":"world"}`))
	check(err)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {

		messageType, p, err := conn.ReadMessage() // read in a message
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("New Message: ", string(p)) // print out that message for clarity
		// errs := conn.WriteMessage(messageType, p)
		for i := 1; i <= len(clients); i++ {
			err = clients[i-1].WriteMessage(messageType, []byte("{'message':'Hello world!'}"))
			check(err)
		}

	}
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	fmt.Println("Hello World")
	http.HandleFunc("/ws", wsEndpoint)
	log.Fatal(http.ListenAndServe(":8765", nil))
}
