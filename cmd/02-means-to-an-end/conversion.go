package main

import (
	"bytes"
	"encoding/binary"
	"log"
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
