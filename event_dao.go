package dao

import (
	"github.com/go-pg/pg"

	"github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

type eventDAO struct {
	DB *pg.DB
}

type event struct {
	model.Event

	tableName struct{} `sql:"select:events_full"`
}

// GetAll returns all events
func (o *eventDAO) GetAll() ([]core.Model, error) {
	events := make([]*event, 0)

	if err := o.DB.Model(&events).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(events))
	for i, e := range events {
		result[i] = core.Model(e)
	}

	return result, nil
}

// GetOne returns one event
func (o *eventDAO) GetOne(id string) (core.Model, error) {
	var e event
	var err error
	e.ID = id
	if err := o.DB.Model(&e).First(); err != nil {
		return nil, err
	}

	return &e, err
}

// Upsert upserts an event
func (o *eventDAO) Upsert(m core.Model) error {
	return o.DB.Insert(m)
}
