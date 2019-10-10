package postgres

import (
	"log"

	"github.com/go-pg/pg"
)

type Logger struct {}

func (Logger) BeforeQuery(q *pg.QueryEvent) {
	log.Println(q.FormattedQuery())
}

func (Logger) AfterQuery(q *pg.QueryEvent) {}