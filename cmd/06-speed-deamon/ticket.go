package main

type Ticket struct {
	plate      string
	road       uint16
	mile       uint16
	timestamp  uint32
	mile2      uint16
	timestamp2 uint32
	speed      uint16
	isSent     bool
}

func NewTicket(plate string, road uint16, mile uint16, timestamp uint32, mile2 uint16, timestamp2 uint32, speed uint16) Ticket {
	return Ticket{plate: plate, road: road, mile: mile, timestamp: timestamp, mile2: mile2, timestamp2: timestamp2, speed: speed}
}
