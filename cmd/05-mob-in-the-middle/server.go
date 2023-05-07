package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
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

	log.Print("Received connection from: " + conn.RemoteAddr().String())

	proxyConn, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		println("Dial failed:", err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s.startProxy(conn, proxyConn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.startProxy(proxyConn, conn)
		wg.Done()
	}()

	wg.Wait()
}

func (s *Server) startProxy(connFrom net.Conn, connTo net.Conn) {
	defer connFrom.Close()
	defer connTo.Close()

	for {
		line, err := s.readLine(connFrom)
		if err != nil {
			return
		}

		line = s.rewriteBogusCoinAddress(line)

		s.writeLine(connTo, line)
		if err != nil {
			return
		}
	}
}

func (s *Server) rewriteBogusCoinAddress(line string) string {
	words := strings.Split(line, " ")
	if len(words) > 0 {
		for index, word := range words {
			if len(word) >= 26 && len(word) <= 35 && word[0] == '7' {
				words[index] = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
			}
		}
		return strings.Join(words, " ")
	} else {
		return line
	}

}

func (s *Server) writeLine(conn net.Conn, line string) error {
	_, err := conn.Write([]byte(line + "\n"))
	return err
}

func (s *Server) readLine(conn net.Conn) (string, error) {
	// Switched to reader because scanner will return the last non-empty line of input if it has no newline
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	return strings.TrimSpace(string(line)), err
}
