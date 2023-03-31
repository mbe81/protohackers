package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	connections map[string]Connection
	ch          chan Message
}

func NewServer() Server {
	return Server{connections: make(map[string]Connection), ch: make(chan Message)}
}

func (s *Server) Run(port int) {

	go s.runMessageDispatcher()

	s.runListener(port)
}

func (s *Server) runMessageDispatcher() {
	for {
		msg := <-s.ch
		for _, s := range s.connections {
			if msg.UserName != s.UserName {
				s.ch <- msg
			}
		}
	}
}

func (s *Server) runListener(port int) {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Accepting TCP connections on port %v", port)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Print("Received connection from: " + conn.RemoteAddr().String())

	// initialize connection
	connection, err := NewConnection(conn)
	if err != nil {
		return
	}
	s.connections[connection.UserName] = connection

	// display online users
	if len(s.connections) > 1 {
		var users []string
		for userName := range s.connections {
			if userName != connection.UserName {
				users = append(users, userName)
			}
		}
		err = connection.writeLine("* Online users: " + strings.Join(users, ", "))
	} else {
		err = connection.writeLine("* You are currently the only user online. 🎉")
	}
	if err != nil {
		return
	}

	// broadcast new user
	s.ch <- NewEvent(connection.UserName+" has entered the room.", connection.UserName)

	// start connection
	connection.Run(s.ch)

	// connection finished
	s.ch <- NewEvent(connection.UserName+" left the building.", connection.UserName)
	delete(s.connections, connection.UserName)

}
