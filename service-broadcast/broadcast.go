package broadcast

import (
	"fmt"
	"log"
	"net"
)

// GetOutboundIP return loval ip
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// StartBroadcast will start service that will send UDP packets on broadcast telling to others
// dispostives what is the IP of socket
func StartBroadcast(broadcastIP string, port int64) {
	var temp string

	temp = fmt.Sprintf(":%d", port)
	pc, err := net.ListenPacket("udp4", temp)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	temp = fmt.Sprintf("%s:%d", broadcastIP, port)
	addr, err := net.ResolveUDPAddr("udp4", temp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", GetOutboundIP())
	temp = fmt.Sprintf()
	_, err := pc.WriteTo([]byte("data to transmit"), addr)
	if err != nil {
		panic(err)
	}
}
