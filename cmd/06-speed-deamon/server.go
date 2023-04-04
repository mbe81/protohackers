package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
}

func NewServer() Server {
	return Server{}
}

func (s *Server) Run(port int) {

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

	fmt.Println("New connection")
	r := bufio.NewReader(conn)
	w := io.Writer(conn)

	var isCamera = false
	var isDispatcher = false
	var hasHeartbeat = false

	for {

		msgType, err := ReadMessageType(r)
		if err != nil {
			WriteErrorMessage(w, "Invalid message type")
			return
		}
		switch msgType {
		case PlateMessageType:
			fmt.Println("Plate")
			msg, err := ReadPlateMessage(r)
			if err != nil {
				WriteErrorMessage(w, "Invalid plate message")
				return
			}
			if isCamera == false {
				WriteErrorMessage(w, "Message not allowed")
				return
			}

			fmt.Println(msg)

		case WantedHeartBeatMessageType:
			fmt.Println("WantHeartbeat")
			msg, err := ReadWantHeartbeatMessages(r)
			if err != nil {
				WriteErrorMessage(w, "Invalid heartbeat message")
				return
			}
			if hasHeartbeat == false {
				h := NewHeartbeat(w)
				h.SetInterval(msg.interval)
				hasHeartbeat = true
			} else {
				WriteErrorMessage(w, "Heartbeat already running")
				return
			}

		case IAmCameraMessageType:
			fmt.Println("IAmCamera")
			msg, err := ReadIAmCameraMessage(r)
			if err != nil {
				WriteErrorMessage(w, "Invalid message type")
				return
			}
			if isCamera == true {
				WriteErrorMessage(w, "Already identified as camera")
				return
			}
			if isDispatcher == true {
				WriteErrorMessage(w, "Already identified as dispatcher")
				return
			}
			isCamera = true

			fmt.Println(msg)

		case IAmDispatcherMessageType:
			fmt.Println("IAmDispatcher")
			msg, err := ReadIAmDispatcherMessage(r)
			if err != nil {
				WriteErrorMessage(w, "Invalid message type")
				return
			}
			if isCamera == true {
				WriteErrorMessage(w, "Already identified as camera")
				return
			}
			if isDispatcher == true {
				WriteErrorMessage(w, "Already identified as dispatcher")
				return
			}
			isDispatcher = true

			fmt.Println(msg)

		}
	}

}

func (s *Server) readLine(conn net.Conn) (string, error) {
	// Switched to reader because scanner will return the last non-empty line of input if it has no newline
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	return strings.TrimSpace(line), err
}
