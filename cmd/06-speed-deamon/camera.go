package main

import "io"

type IAmCameraMessage struct {
	road  uint16
	mile  uint16
	limit uint16
}

func ReadIAmCameraMessage(r io.Reader) (IAmCameraMessage, error) {
	var err error

	msg := IAmCameraMessage{}
	if msg.road, err = ReadUInt16(r); err != nil {
		return msg, err
	}
	if msg.mile, err = ReadUInt16(r); err != nil {
		return msg, err
	}
	if msg.limit, err = ReadUInt16(r); err != nil {
		return msg, err
	}
	return msg, err
}
