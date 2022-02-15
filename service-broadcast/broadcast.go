package broadcast

import (
	"fmt"
	"log"
	"net"
	"time"
)

// GetOutboundIP return loval ip
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// StartBroadcast will start service that will send UDP packets on broadcast telling to others
// dispostives what is the IP of socket
func StartBroadcast(broadcastIP string, portUDP int64, portWeb int64, interval time.Duration) {
	temp := fmt.Sprintf(":%d", portUDP)
	pc, err := net.ListenPacket("udp4", temp)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		fmt.Printf("%s sent this: %s\n", addr, buf[:n])

		temp = fmt.Sprintf("%s:%d\r\n", GetOutboundIP(), portWeb)

		go pc.WriteTo([]byte(temp), addr)

	}

}
