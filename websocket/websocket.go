package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var dispatcherFunc = func([]byte) WebSocketPacket { return WebSocketPacket{} }

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%d Bytes Received: %s\n", len(p), string(p))
		commandBack := dispatcherFunc(p)
		response, _ := json.Marshal(commandBack)
		conn.WriteMessage(messageType, response)
		log.Printf("%d bytes sent\n", len(response))
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	reader(ws)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func setupRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", handle)
}

// StartServer WebServer
func StartServer(param string, dispatcher func(msg []byte) WebSocketPacket) {
	dispatcherFunc = dispatcher
	log.Println("Starting Server")
	setupRoutes()
	log.Fatal(http.ListenAndServe(param, nil))
}
