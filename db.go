package dao

import (
	"log"

	"github.com/go-pg/pg"
)

// Init establishes the database connection and returns the client
func Init(o *pg.Options) *pg.DB {
	log.Println("Initializing the DB with options... | ", o)
	if o.Addr == "" {
		o.Addr = "localhost:5432"
	}
	return pg.Connect(o)
}
