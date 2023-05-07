package main

import (
	"io"
	"time"
)

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
