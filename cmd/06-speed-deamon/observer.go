package main

import (
	"fmt"
	"sync"
	"time"
)

type Plate struct {
	road      uint16
	mile      uint16
	limit     uint16
	plate     string
	timestamp uint32
}

func NewPlate(road uint16, mile uint16, limit uint16, plate string, timestamp uint32) Plate {
	return Plate{road: road, mile: mile, limit: limit, plate: plate, timestamp: timestamp}
}

type Ticket struct {
	plate      string
	road       uint16
	mile       uint16
	timestamp  uint32
	mile2      uint16
	timestamp2 uint32
	speed      uint16
	isReserved bool
}

func NewTicket(plate string, road uint16, mile uint16, timestamp uint32, mile2 uint16, timestamp2 uint32, speed uint16) Ticket {
	return Ticket{plate: plate, road: road, mile: mile, timestamp: timestamp, mile2: mile2, timestamp2: timestamp2, speed: speed}
}

type Dispatcher struct {
	roads []uint16
}

type Observer struct {
	mux     sync.Mutex
	plates  []Plate
	tickets []Ticket
}

func NewObserver() Observer {
	var o = Observer{}
	return o
}

func (o *Observer) NewPlate(p Plate) {
	o.mux.Lock()
	defer o.mux.Unlock()

	fmt.Println("NEW Plate", p)
	fmt.Println("Current Plates", o.plates)
	o.plates = append(o.plates, p)
	fmt.Println("New Plates", o.plates)
}

func (o *Observer) RemovePlate(p Plate) {
	o.mux.Lock()
	defer o.mux.Unlock()

	for i, plate := range o.plates {

		if plate.road == p.road &&
			plate.mile == p.mile &&
			plate.plate == p.plate &&
			plate.timestamp == p.timestamp {

			o.plates[i] = o.plates[len(o.plates)-1]
			o.plates = o.plates[:len(o.plates)-1]
			break
		}
	}
}

func (o *Observer) NewTicket(t Ticket) {
	o.mux.Lock()
	defer o.mux.Unlock()

	o.tickets = append(o.tickets, t)

}

func (o *Observer) Observe() {
	for {
		time.Sleep(2 * time.Second)

		fmt.Println("Current Tickets", o.tickets)

		for _, plate1 := range o.plates {
			for _, plate2 := range o.plates {

				if plate1.road == plate2.road &&
					plate1.plate == plate2.plate &&
					plate1.timestamp < plate2.timestamp {

					fmt.Println("Valid Observation", plate1, plate2)

					t := plate2.timestamp - plate1.timestamp
					d := plate2.mile - plate1.mile

					s := uint16(float32(d) / float32(t) * 3600 * 100)
					fmt.Println("t", t, "d", d, "s", s)

					if s > plate2.limit {
						ticket := NewTicket(plate1.plate,
							plate1.road,
							plate1.mile,
							plate1.timestamp,
							plate2.mile,
							plate2.timestamp,
							s)

						fmt.Println("Ticket!!!!!", ticket)
						//o.RemovePlate(plate1)
						//o.RemovePlate(plate2)
						o.NewTicket(ticket)
					}

				}

			}
		}

	}
}

func (o *Observer) GetTicket(road uint16) *Ticket {
	o.mux.Lock()
	defer o.mux.Unlock()

	for i, ticket := range o.tickets {
		if ticket.road == road &&
			ticket.isReserved == false {
			o.tickets[i].isReserved = true
			return &ticket
		}
	}

	return nil
}

func (o *Observer) RemoveTicket(t Ticket) {
	o.mux.Lock()
	defer o.mux.Unlock()

	for i, ticket := range o.tickets {
		if ticket.road == t.road &&
			ticket.mile == t.mile &&
			t.plate == t.plate {

			o.tickets[i] = o.tickets[len(o.tickets)-1]
			o.tickets = o.tickets[:len(o.tickets)-1]
		}
	}

}
