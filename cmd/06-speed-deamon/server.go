package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	mux      sync.Mutex
	cPlates  chan Event
	plates   []Event
	cTickets chan Ticket
	tickets  []Ticket
}

func NewServer() *Server {
	s := Server{}
	s.cPlates = make(chan Event)
	s.cTickets = make(chan Ticket)
	return &s
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

	go s.HandlePlates()
	go s.HandleTickets()

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
			//fmt.Println("IAmCamera")
			camera, err := ReadIAmCameraMessage(r)
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
			currentRoad = camera.road
			currentMile = camera.mile
			currentLimit = camera.limit

		case PlateMessageType:
			//fmt.Println("Receive plate")
			plateMessage, err := ReadPlateMessage(r)
			if err != nil {
				WriteErrorMessage(w, "Invalid plate message")
				return
			}
			if isCamera == false {
				WriteErrorMessage(w, "Message not allowed")
				return
			}

			var plate = Event{}
			plate.plate = plateMessage.plate
			plate.timestamp = plateMessage.timestamp
			plate.road = currentRoad
			plate.mile = currentMile
			plate.limit = currentLimit

			s.cPlates <- plate

		case WantedHeartBeatMessageType:
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
					_ = s.GetTicket(w, road)
					//if t != nil {
					//	_ = t.Write(w)
					//}
				}
				time.Sleep(time.Millisecond * 100)
			}

		}
	}

}

func (s *Server) HandlePlates() {

	var duration uint32
	var distance, speed uint16
	//var index int
	var oldPlate Event
	var newTicket Ticket

	for {
		newPlate := <-s.cPlates

		s.mux.Lock()

		//fmt.Println("New plate", newPlate)
		//violation := false
		for _, oldPlate = range s.plates {

			if newPlate.road == oldPlate.road && newPlate.plate == oldPlate.plate {

				if newPlate.timestamp < oldPlate.timestamp {
					tempPlate := newPlate
					newPlate = oldPlate
					oldPlate = tempPlate
				}

				duration = newPlate.timestamp - oldPlate.timestamp
				if newPlate.mile > oldPlate.mile {
					distance = newPlate.mile - oldPlate.mile
				} else {
					distance = oldPlate.mile - newPlate.mile
				}

				if duration > 0 {
					speed = uint16(float32(distance) / float32(duration) * 3600 * 100)
				}

				if speed > newPlate.limit*100 {
					//violation = true
					newTicket = NewTicket(newPlate.plate,
						newPlate.road,
						oldPlate.mile,
						oldPlate.timestamp,
						newPlate.mile,
						newPlate.timestamp,
						speed)

					s.cTickets <- newTicket
					fmt.Println("New ticket", newTicket)
					break

				}
			}
		}

		s.plates = append(s.plates, newPlate)
		s.mux.Unlock()

	}
}

func (s *Server) HandleTickets() {

	var oldTicket Ticket

	for {
		newTicket := <-s.cTickets
		s.mux.Lock()
		saveTicket := true

		for _, oldTicket = range s.tickets {
			if newTicket.plate == oldTicket.plate {
				// ticket for same plate and road already exists
				// if it is on the same day do not save it
				if newTicket.timestamp/86400 == oldTicket.timestamp2/86400 {
					saveTicket = false
				}
				//if newTicket.timestamp2/86400 == oldTicket.timestamp/86400 {
				//	saveTicket = false
				//}
				//if newTicket.timestamp2/86400 == oldTicket.timestamp2/86400 {
				//	saveTicket = false
				//}
			}
		}

		if saveTicket {
			s.tickets = append(s.tickets, newTicket)
		}
		s.mux.Unlock()

	}
}

func (s *Server) GetTicket(w io.Writer, road uint16) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for index := 0; index < len(s.tickets); index++ {
		if s.tickets[index].road == road && s.tickets[index].isSent == false {
			err := WriteTicket(w, s.tickets[index])
			if err != nil {
				return err
			}
			s.tickets[index].isSent = true
		}
	}

	return nil
}
