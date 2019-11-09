package postgres

import (
	"context"
	"log"

	"github.com/go-pg/pg"
)

// Logger is the logger plugin for pg
type Logger struct{}

// BeforeQuery implements pg.QueryHook
func (Logger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	log.Println(q.FormattedQuery())
	return c, nil
}

// AfterQuery implements pg.QueryHook
func (Logger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	return nil
}
