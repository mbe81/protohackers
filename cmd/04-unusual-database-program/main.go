package main

import (
	"fmt"
	"net"
	"protohackers/internal/udp"
	"strings"
)

var data = make(map[string]string)

func handler(conn *net.UDPConn) {

	for {
		message := make([]byte, 255)
		rlen, addr, err := conn.ReadFromUDP(message[:])
		if err != nil {
			panic(err)
		}

		queryParts := strings.SplitN(string(message[:rlen]), "=", 2)

		if queryParts[0] == "version" {
			response := "version=Key Value Store v1.0"
			conn.WriteToUDP([]byte(response), addr)
		} else if len(queryParts) == 2 {
			// insert
			fmt.Println("Insert", string(message[:rlen]), addr)
			key := queryParts[0]
			value := queryParts[1]
			data[key] = value

		} else {
			// query
			fmt.Println("Query", string(message[:rlen]), addr)
			key := queryParts[0]
			value := data[key]
			response := key + "=" + value
			conn.WriteToUDP([]byte(response), addr)
		}
	}
}

func main() {
	udp.RunUDPServer(handler, 9001)
}
