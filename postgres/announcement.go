package postgres

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// AnnouncementQueryStrategy implements QueryStrategy for announcements
type AnnouncementQueryStrategy struct {
	*pg.DB
}

type announcement struct {
	model.Announcement

	tableName struct{} `sql:"select:announcements_full"`
}

// ModelName outputs this model's name
func (s *AnnouncementQueryStrategy) ModelName() string {
	return "announcement"
}

// GetAll returns all announcements
func (s *AnnouncementQueryStrategy) GetAll() ([]core.Model, error) {
	as := make([]*announcement, 0)

	if err := s.DB.Model(&as).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(as))
	for i, a := range as {
		result[i] = core.Model(a)
	}

	return result, nil
}

// GetOne returns one announcement
func (s *AnnouncementQueryStrategy) GetOne(id string) (core.Model, error) {
	var a announcement
	a.ID = id
	if err := s.DB.Model(&a).First(); err != nil {
		return nil, err
	}

	return &a, nil
}

// Search finds all announcements meeting the criteria given by the payload
func (s *AnnouncementQueryStrategy) Search(payload []byte) ([]core.Model, error) {
	var params announcement
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	as := make([]*announcement, 0)

	query := s.DB.Model(&as)

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
func (s *AnnouncementQueryStrategy) Upsert(m core.Model) error {
	_, err := s.DB.Model(m).
		OnConflict("(id) DO UPDATE").
		Set(`(
			name,
			content
		) = (
			?name,
			?content)`).
		Insert(m)
	return err
}
