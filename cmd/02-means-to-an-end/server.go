package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
}

func NewServer() Server {
	var s Server
	return s
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

	log.Print("Received connection from: " + conn.RemoteAddr().String())

	// Each session can only query the data supplied by itself so the prices is defined in the connection handler
	var (
		db     = NewDatabase()
		buffer = make([]byte, 9)
		result int32
	)

	for {
		if _, err := io.ReadFull(conn, buffer); err != nil {
			log.Print("Error reading data from: " + conn.RemoteAddr().String())
			return
		}

		action := string(buffer[:1])
		result = -1

		if action == "I" {
			timestamp := byteToInt32(buffer[1:5])
			price := byteToInt32(buffer[5:9])
			db.InsertPrice(timestamp, price)
		} else if action == "Q" {
			minTime := byteToInt32(buffer[1:5])
			maxTime := byteToInt32(buffer[5:9])
			result = db.RequestAverage(minTime, maxTime)
		}

		if result != -1 {
			if _, err := conn.Write(int32ToByte(result)); err != nil {
				log.Print("Error writing data to: " + conn.RemoteAddr().String())
				return
			}
		}
	}
}
