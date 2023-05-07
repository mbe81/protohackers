package main

type PlateHistory struct {
	events  []Event
	tickets []Ticket
}
type Event struct {
	road      uint16
	mile      uint16
	limit     uint16
	plate     string
	timestamp uint32
}
