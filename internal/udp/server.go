package udp

import (
	"fmt"
	"log"
	"net"
)

type Handler func(conn *net.UDPConn)

func RunUDPServer(handler Handler, port int) {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Accepting UDP connections on port %v", port)

	handler(listener)

}
