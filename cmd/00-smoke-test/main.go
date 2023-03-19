package main

import (
	"io"
	"log"
	"net"
	"protohackers/internal/tcp"
)

type SmokeTest struct{}

func (SmokeTest) Handle(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}
}

func main() {
	tcp.RunServer(SmokeTest{}, 9001)
}
