package main

import (
	"errors"
	"io"
)

type MessageType uint8

const (
	ErrorMessageType           MessageType = 0x10 // server
	PlateMessageType           MessageType = 0x20 // client
	TicketMessageType          MessageType = 0x21 // server
	WantedHeartBeatMessageType MessageType = 0x40 // client
	HeartbeatMessageType       MessageType = 0x41 // server
	IAmCameraMessageType       MessageType = 0x80 // client
	IAmDispatcherMessageType   MessageType = 0x81 // client
	UnknownMessageType         MessageType = 0xFF
)

func ReadMessageType(r io.Reader) (MessageType, error) {
	b, err := ReadByte(r)
	if err != nil {
		return UnknownMessageType, err
	}

	switch MessageType(b) {
	case PlateMessageType, WantedHeartBeatMessageType, IAmCameraMessageType, IAmDispatcherMessageType:
		return MessageType(b), nil
	}
	return UnknownMessageType, errors.New("unknown message type")
}

func WriteMessageType(w io.Writer, t MessageType) error {
	if err := WriteByte(w, byte(t)); err != nil {
		return err
	}
	return nil
}
