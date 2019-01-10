package dao

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

type announcementDAO struct {
	DB *pg.DB
}

type announcement struct {
	model.Announcement

	tableName struct{} `sql:"select:announcements_full"`
}

// GetAll returns all announcements
func (o *announcementDAO) GetAll() ([]core.Model, error) {
	as := make([]*announcement, 0)

	if err := o.DB.Model(&as).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(as))
	for i, a := range as {
		result[i] = core.Model(a)
	}

	return result, nil
}

// GetOne returns one announcement
func (o *announcementDAO) GetOne(id string) (core.Model, error) {
	var a announcement
	var err error
	a.ID = id
	if err := o.DB.Model(&a).First(); err != nil {
		return nil, err
	}

	return &a, err
}

// Search finds all announcements meeting the criteria given by the payload
func (o *announcementDAO) Search(payload []byte) ([]core.Model, error) {
	var params announcement
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	as := make([]*announcement, 0)

	query := o.DB.Model(&as)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		query = query.Where("name = ?", params.Name)
	}

	if params.AuthorID != "" {
		query = query.Where("author_id = ?", params.AuthorID)
	}

	if err := query.Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(as))
	for i, a := range as {
		result[i] = core.Model(a)
	}

	return result, nil
}

// Upsert upserts an announcement
func (o *announcementDAO) Upsert(m core.Model) error {
	return o.DB.Insert(m)
}
