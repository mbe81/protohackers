package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"protohackers/internal/tcp"
)

func byteToInt32(b []byte) (i int32) {
	buffer := bytes.NewBuffer(b)
	err := binary.Read(buffer, binary.BigEndian, &i)
	if err != nil {
		log.Fatal("Converting bytes to int32 failed")
	}
	return
}

func int32ToByte(i int32) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, i)
	if err != nil {
		log.Fatal("Converting int32 to bytes failed")
	}
	return buffer.Bytes()
}

func processMessage(m []byte, prices map[int32]int32) int32 {
	action := string(m[:1])
	if action == "I" {
		timestamp := byteToInt32(m[1:5])
		price := byteToInt32(m[5:9])
		return processInsert(timestamp, price, prices)
	} else if action == "Q" {
		minTime := byteToInt32(m[1:5])
		maxTime := byteToInt32(m[5:9])
		return processQuery(minTime, maxTime, prices)
	} else {
		return -1
	}

}

func processInsert(timestamp int32, price int32, prices map[int32]int32) int32 {
	prices[timestamp] = price
	return -1
}

func processQuery(minTime int32, maxTime int32, prices map[int32]int32) int32 {
	avgPrice := 0
	if minTime > maxTime {
		avgPrice = 0
	} else {
		sum := 0
		count := 0

		for timestamp, price := range prices {
			if timestamp >= minTime && timestamp <= maxTime {
				sum += int(price)
				count++
			}
		}

		if count > 0 {
			avgPrice = sum / count
		} else {
			avgPrice = 0
		}

	}

	return int32(avgPrice)
}

func main() {
	tcp.RunTCPServer(handler, 9001)
}

func handler(conn net.Conn) {
	defer conn.Close()

	log.Print("Received connection from: " + conn.RemoteAddr().String())

	// Each session can only query the data supplied by itself so map is defined in handler
	var prices = make(map[int32]int32)
	buffer := make([]byte, 9)

	for {
		if _, err := io.ReadFull(conn, buffer); err != nil {
			log.Print("Close connection from: " + conn.RemoteAddr().String())
			conn.Close()
			return
		}

		result := processMessage(buffer, prices)
		if result != -1 {
			conn.Write(int32ToByte(result))
		}
	}
}
