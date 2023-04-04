package main

import (
	"io"
)

type PlateMessage struct {
	plate     string
	timestamp uint32
}

func ReadPlateMessage(r io.Reader) (PlateMessage, error) {
	var err error

	msg := PlateMessage{}
	if msg.plate, err = ReadString(r); err != nil {
		return msg, err
	}
	if msg.timestamp, err = ReadUInt32(r); err != nil {
		return msg, err
	}
	return msg, err
}
