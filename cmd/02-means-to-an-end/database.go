package main

type Database struct {
	prices map[int32]int32
}

func NewDatabase() Database {
	db := Database{}
	db.prices = make(map[int32]int32)
	return db
}

func (db *Database) InsertPrice(timestamp int32, price int32) {
	db.prices[timestamp] = price
}

func (db *Database) RequestAverage(minTime int32, maxTime int32) int32 {
	avgPrice := 0
	if minTime > maxTime {
		avgPrice = 0
	} else {
		sum := 0
		count := 0

		for timestamp, price := range db.prices {
			if timestamp >= minTime && timestamp <= maxTime {
				sum += int(price)
				count++
			}
		}

		if count > 0 {
			avgPrice = sum / count
		} else {
			avgPrice = 0
		}

	}

	return int32(avgPrice)
}
