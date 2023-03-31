package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Server struct {
}

type RequestDTO struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type ResponseDTO struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
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
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		buffer := scanner.Bytes()
		log.Println(string(buffer))

		req := RequestDTO{}
		err := json.Unmarshal(buffer, &req)
		if err != nil || req.Method != "isPrime" || req.Number == nil {
			conn.Write([]byte("Error!"))
			conn.Close()
			return
		}

		res, _ := json.Marshal(ResponseDTO{Method: req.Method, Prime: isPrime(int(*req.Number))})
		conn.Write(res)
		conn.Write([]byte("\n"))
	}
}
