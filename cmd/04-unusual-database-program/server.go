package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	conn *net.UDPConn
	db   Database
}

func NewServer() Server {
	var s Server
	s.db = NewDatabase()
	return s
}

func (s *Server) requestHandler() {
	for {
		message := make([]byte, 255)
		length, addr, err := s.conn.ReadFromUDP(message[:])
		if err != nil {
			fmt.Println("Error reading request")
			continue
		}

		queryParts := strings.SplitN(string(message[:length]), "=", 2)

		if queryParts[0] == "version" {
			// version
			key := queryParts[0]
			value := s.db.Version()
			s.writeResponse(addr, key, value)

		} else if len(queryParts) == 2 {
			// insert
			key := queryParts[0]
			value := queryParts[1]
			s.db.Insert(key, value)

		} else {
			// retrieve
			key := queryParts[0]
			value := s.db.Retrieve(key)
			s.writeResponse(addr, key, value)
			if err != nil {
				fmt.Println("Error writing response")
				continue
			}

		}
	}
}

func (s *Server) writeResponse(addr *net.UDPAddr, key string, value string) error {
	response := key + "=" + value
	_, err := s.conn.WriteToUDP([]byte(response), addr)
	return err
}

func (s *Server) Run(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Accepting UDP connections on port %v", port)

	s.requestHandler()

}
