package main

import (
	"fmt"
	"io"
)

type IAmDispatcherMessage struct {
	numRoads uint8
	roads    []uint16
}

func ReadIAmDispatcherMessage(r io.Reader) (IAmDispatcherMessage, error) {
	var err error
	var numRoads uint8

	msg := IAmDispatcherMessage{}
	if numRoads, err = ReadUInt8(r); err != nil {
		return msg, err
	}

	for i := 1; i <= int(numRoads); i++ {
		var road uint16
		if road, err = ReadUInt16(r); err != nil {
			return msg, err
		}

		msg.roads = append(msg.roads, road)
	}

	msg.numRoads = numRoads
	fmt.Println(msg)
	return msg, err
}
