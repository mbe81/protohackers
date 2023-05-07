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
	return msg, err
}

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

func WriteTicket(w io.Writer, t Ticket) error {
	if err := WriteMessageType(w, TicketMessageType); err != nil {
		return err
	}
	if err := WriteString(w, t.plate); err != nil {
		return err
	}
	if err := WriteUint16(w, t.road); err != nil {
		return err
	}
	if err := WriteUint16(w, t.mile); err != nil {
		return err
	}
	if err := WriteUint32(w, t.timestamp); err != nil {
		return err
	}
	if err := WriteUint16(w, t.mile2); err != nil {
		return err
	}
	if err := WriteUint32(w, t.timestamp2); err != nil {
		return err
	}
	if err := WriteUint16(w, t.speed); err != nil {
		return err
	}
	return nil
}

func WriteErrorMessage(w io.Writer, msg string) {
	if err := WriteMessageType(w, ErrorMessageType); err != nil {
		return
	}
	if err := WriteString(w, msg); err != nil {
		return
	}
	return
}

type WantHeartbeatMessage struct {
	interval uint32
}

func ReadWantHeartbeatMessages(r io.Reader) (WantHeartbeatMessage, error) {
	var err error
	msg := WantHeartbeatMessage{}
	if msg.interval, err = ReadUInt32(r); err != nil {
		return msg, err
	}
	return msg, err
}

func WriteHeartbeatMessage(w io.Writer) error {
	// Heartbeat message is an empty message. Only the message type is written
	if err := WriteMessageType(w, HeartbeatMessageType); err != nil {
		return err
	}
	return nil
}
