package dao

import (
	"github.com/go-pg/pg"

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

// Upsert upserts an announcement
func (o *announcementDAO) Upsert(m core.Model) error {
	return o.DB.Insert(m)
}
