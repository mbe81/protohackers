package main

import (
	"io"
	"log"
	"net"
	"protohackers/internal/tcp"
)

func handler(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}
}

func main() {
	tcp.RunTCPServer(handler, 9001)
}
