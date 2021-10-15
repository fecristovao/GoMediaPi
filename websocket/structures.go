package websocket

type WebSocketPacket struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}
