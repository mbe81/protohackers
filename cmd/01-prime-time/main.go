package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"protohackers/internal/tcp"
)

type PrimeTime struct{}

type RequestDTO struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type ResponseDTO struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func (PrimeTime) Handle(conn net.Conn) {
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

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	tcp.RunServer(PrimeTime{}, 9001)
}
