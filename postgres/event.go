package postgres

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// EventQueryStrategy implements QueryStrategy for events
type EventQueryStrategy struct {
	*pg.DB
}

type event struct {
	model.Event

	tableName struct{} `sql:"select:events_full"`
}

// ModelName outputs this model's name
func (s *EventQueryStrategy) ModelName() string {
	return "event"
}

// GetAll returns all events
func (s *EventQueryStrategy) GetAll() ([]core.Model, error) {
	events := make([]*event, 0)

	if err := s.DB.Model(&events).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(events))
	for i, e := range events {
		result[i] = core.Model(e)
	}

	return result, nil
}

// GetOne returns one event
func (s *EventQueryStrategy) GetOne(id string) (core.Model, error) {
	var e event
	e.ID = id
	if err := s.DB.Model(&e).WherePK().First(); err != nil {
		return nil, err
	}

	return &e, nil
}

// Search finds all events meeting the criteria given by the payload
func (s *EventQueryStrategy) Search(payload []byte) ([]core.Model, error) {
	var params event
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	es := make([]*event, 0)

	query := s.DB.Model(&es)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		query = query.Where("name ilike ?", "%"+params.Name+"%")
	}

	if params.Date != nil {
		query = query.Where("date = ?", *params.Date)
	}

	if err := query.Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(es))
	for i, e := range es {
		result[i] = core.Model(e)
	}

	return result, nil
}

// Upsert upserts an event
func (s *EventQueryStrategy) Upsert(m core.Model) error {
	_, err := s.DB.Model(m).
		OnConflict("(id) DO UPDATE").
		Set(`(
			name,
			length,
			deadline,
			allow_maybe,
			description,
			location,
			location_description
		) = (
			?name,
			?length,
			?deadline,
			?allow_maybe,
			?description,
			?location,
			?location_description)`).
		Insert(m)
	return err
}

// UpsertRelationship upserts an event relationship
func (s *EventQueryStrategy) UpsertRelationship(e core.Entity, relation string) error {
	return nil
}
