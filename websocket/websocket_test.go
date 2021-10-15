package websocket

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/fecristovao/GoModPi/m3u"
)

func TestServer(t *testing.T) {
	StartServer(":8090", func(msg []byte) WebSocketPacket {
		packet := WebSocketPacket{}

		if err := json.Unmarshal(msg, &packet); err != nil {
			log.Println(err)
			return packet
		}

		switch packet.Command {
		case "getAllChannels":
			packet.Data = m3u.ParseFile("../m3u/iptv.m3u")
		}

		return packet
	})
}
