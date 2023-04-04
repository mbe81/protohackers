package main

import (
	"io"
	"time"
)

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

type Heartbeat struct {
	currentInterval uint32
	w               io.Writer
}

func NewHeartbeat(w io.Writer) Heartbeat {
	return Heartbeat{0, w}
}

func (h *Heartbeat) SetInterval(interval uint32) {
	h.currentInterval = interval
	if h.currentInterval > 0 {
		go h.Beat(h.w, h.currentInterval)
	}
}

func (h *Heartbeat) Beat(w io.Writer, interval uint32) {
	if interval == h.currentInterval {
		if err := WriteHeartbeatMessage(w); err != nil {
			return
		}
		time.Sleep(time.Duration(h.currentInterval*100) * time.Millisecond)
		h.Beat(w, interval)
	}
}

func (h *Heartbeat) IsRunning() bool {
	return h.currentInterval > 0
}
