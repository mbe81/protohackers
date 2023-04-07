package main

import "io"

type TicketMessage struct {
	plate      string
	road       uint16
	mile       uint16
	timestamp  uint32
	mile2      uint16
	timestamp2 uint32
	speed      uint16
}

func WriteTicketMessage(w io.Writer, msg TicketMessage) error {
	if err := WriteMessageType(w, TicketMessageType); err != nil {
		return err
	}
	if err := WriteString(w, msg.plate); err != nil {
		return err
	}
	if err := WriteUint16(w, msg.road); err != nil {
		return err
	}
	if err := WriteUint16(w, msg.mile); err != nil {
		return err
	}
	if err := WriteUint32(w, msg.timestamp); err != nil {
		return err
	}
	if err := WriteUint16(w, msg.mile2); err != nil {
		return err
	}
	if err := WriteUint32(w, msg.timestamp2); err != nil {
		return err
	}
	if err := WriteUint16(w, msg.speed); err != nil {
		return err
	}
	return nil
}
