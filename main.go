package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/fecristovao/GoModPi/m3u"
	"github.com/fecristovao/GoModPi/vlc"
	"github.com/fecristovao/GoModPi/websocket"
)

var parsedIPTV = m3u.ParseFile("/home/pi/GoMediaPi/iptv.m3u")
var vlcPath = "/usr/bin/vlc"
var vlcPID = -1

var commands = map[string]func(interface{}) websocket.WebSocketPacket{
	"getAllChannels": getAllChannels,
	"watchChannel":   watchChannel,
	"stop":           stopVLC,
}

func stopVLC(interface{}) websocket.WebSocketPacket {
	vlc.CloseVLC(vlcPID)
	return websocket.WebSocketPacket{
		Command: "stop",
		Data:    "",
	}
}

func watchChannel(channel interface{}) websocket.WebSocketPacket {
	response := websocket.WebSocketPacket{
		Command: "watchChannel",
	}
	selectedChannel := channel.(map[string]interface{})
	response.Data = selectedChannel

	link := strings.Trim(selectedChannel["StreamLink"].(string), "\r")
	link = strings.Trim(link, "\r\n")
	vlcPID, _ = vlc.OpenVLC(vlcPath, "--fullscreen --no-qt-error-dialogs --one-instance", link)

	return response
}

func getAllChannels(interface{}) websocket.WebSocketPacket {
	response := websocket.WebSocketPacket{
		Command: "getAllChannels",
		Data:    parsedIPTV,
	}

	return response
}

func dispatcher(msg []byte) websocket.WebSocketPacket {
	packet := websocket.WebSocketPacket{}

	if err := json.Unmarshal(msg, &packet); err != nil {
		log.Println(err)
		return packet
	}

	for command, action := range commands {
		if command == packet.Command {
			packet = action(packet.Data)
			break
		}
	}

	return packet
}

func main() {
	websocket.StartServer(":8090", dispatcher)
}
