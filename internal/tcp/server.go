package tcp

import (
	"fmt"
	"log"
	"net"
)

type Server interface {
	Handle(conn net.Conn)
}

func RunServer(s Server, port int) {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Accepting TCP connections on port %v", port)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.Handle(conn)
	}
}
