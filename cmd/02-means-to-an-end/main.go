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

func processInsert(timestamp int32, price int32, prices map[int32]int32) {
	prices[timestamp] = price
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

func handler(conn net.Conn) {
	defer conn.Close()

	log.Print("Received connection from: " + conn.RemoteAddr().String())

	// Each session can only query the data supplied by itself so the prices is defined in the connection handler
	var (
		prices = make(map[int32]int32)
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
			processInsert(timestamp, price, prices)
		} else if action == "Q" {
			minTime := byteToInt32(buffer[1:5])
			maxTime := byteToInt32(buffer[5:9])
			result = processQuery(minTime, maxTime, prices)
		}

		if result != -1 {
			if _, err := conn.Write(int32ToByte(result)); err != nil {
				log.Print("Error writing data to: " + conn.RemoteAddr().String())
				return
			}
		}
	}
}

func main() {
	tcp.RunTCPServer(handler, 9001)
}
