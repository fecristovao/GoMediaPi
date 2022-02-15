package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	broadcast "github.com/fecristovao/GoModPi/service-broadcast"
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
	go reader(ws)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func discover(w http.ResponseWriter, r *http.Request) {
	tmp := fmt.Sprintf("%s:8090\n", broadcast.GetOutboundIP())
	fmt.Fprintf(w, tmp)
}

func setupRoutes() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handle)
	http.HandleFunc("/discover", discover)
}

// StartServer WebServer
func StartServer(param string, dispatcher func(msg []byte) WebSocketPacket) {
	dispatcherFunc = dispatcher
	log.Println("Starting Server")
	setupRoutes()
	log.Fatal(http.ListenAndServe(param, nil))
}
