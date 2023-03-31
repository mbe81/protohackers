package main

type Database struct {
	data map[string]string
}

func NewDatabase() Database {
	var db Database
	db.data = make(map[string]string)
	return db
}

func (db *Database) Insert(key string, value string) {
	db.data[key] = value
}

func (db *Database) Retrieve(key string) string {
	return db.data[key]
}

func (db *Database) Version() string {
	return "UDP - Universal Database Program v1.0"
}
