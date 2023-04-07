package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	observer Observer
}

func NewServer() Server {
	return Server{observer: NewObserver()}
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

	go s.observer.Observe()

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

	var (
		isCamera     = false
		isDispatcher = false
		hasHeartbeat = false

		currentRoad  = uint16(0)
		currentMile  = uint16(0)
		currentLimit = uint16(0)
	)

	for {

		msgType, err := ReadMessageType(r)
		if err != nil {
			WriteErrorMessage(w, "Invalid message type")
			return
		}
		switch msgType {
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
			currentRoad = msg.road
			currentMile = msg.mile
			currentLimit = msg.limit

			fmt.Println(msg)

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

			s.observer.NewPlate(NewPlate(currentRoad, currentMile, currentLimit, msg.plate, msg.timestamp))

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

			for {
				for _, road := range msg.roads {
					fmt.Println("GET TICKETS FOR ROAD", road)
					t := s.observer.GetTicket(road)
					if t != nil {
						fmt.Println("Ticket to DISPATCH", t)
						err := WriteTicketMessage(w, TicketMessage{
							plate:      t.plate,
							road:       t.road,
							mile:       t.mile,
							timestamp:  t.timestamp,
							mile2:      t.mile2,
							timestamp2: t.timestamp2,
							speed:      t.speed,
						})
						if err != nil {
							s.observer.RemoveTicket(*t)
						}
						s.observer.RemoveTicket(*t)
					}

				}

				time.Sleep(time.Second)
			}

		}
	}

}
