package main

import (
	"log"
	"net"
	"protohackers/internal/tcp"
	"strings"
)

type Server struct {
	sessions map[string]Session
	channel  chan Message
}

func NewServer() Server {
	return Server{sessions: make(map[string]Session), channel: make(chan Message)}
}

func (s *Server) Run() {
	go s.StartDispatcher()
	s.StartListener()
}

func (s *Server) StartDispatcher() {
	for {
		msg := <-s.channel
		for _, s := range s.sessions {
			if msg.UserName != s.UserName {
				s.channel <- msg
			}
		}
	}
}

func (s *Server) StartListener() {
	tcp.RunTCPServer(s.SessionHandler, 9001)
}

func (s *Server) SessionHandler(conn net.Conn) {
	defer conn.Close()

	log.Print("Received connection from: " + conn.RemoteAddr().String())

	// initialize session
	session, err := NewSession(conn)
	if err != nil {
		return
	}
	s.sessions[session.UserName] = session

	// display online users
	if len(s.sessions) > 1 {
		var users []string
		for userName := range s.sessions {
			if userName != session.UserName {
				users = append(users, userName)
			}
		}
		err = session.writeLine("* Online users: " + strings.Join(users, ", "))
	} else {
		err = session.writeLine("* You are currently the only user online. 🎉")
	}
	if err != nil {
		return
	}

	// broadcast new user
	s.channel <- NewEvent(session.UserName+" has entered the room.", session.UserName)

	// start session
	session.Run(s.channel)

	// session finished
	s.channel <- NewEvent(session.UserName+" left the building.", session.UserName)
	delete(s.sessions, session.UserName)

}
