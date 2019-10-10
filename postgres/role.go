package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/v9/orm"
	"github.com/pkg/errors"

	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// RoleQueryStrategy implements QueryStrategy for roles
type RoleQueryStrategy struct {
	*pg.DB
}

type role struct {
	model.Role
}

// ModelName outputs this model's name
func (s *RoleQueryStrategy) ModelName() string {
	return "role"
}

// GetAll returns all roles
func (s *RoleQueryStrategy) GetAll() ([]core.Model, error) {
	roles := make([]*role, 0)

	err := s.DB.Model(&roles).Select()
	if err != nil {
		return nil, err
	}

	result := make([]core.Model, len(roles))
	for i, r := range roles {
		result[i] = core.Model(r)
	}

	return result, nil
}

// GetOne returns one role
func (s *RoleQueryStrategy) GetOne(id string) (core.Model, error) {
	var r role
	r.ID = id
	if err := s.DB.Model(&r).WherePK().First(); err != nil {
		return nil, err
	}

	return &r, nil
}

// Search finds all roles meeting the criteria given by the payload
func (s *RoleQueryStrategy) Search(payload []byte) ([]core.Model, error) {
	var params role
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	rs := make([]*role, 0)

	query := s.DB.Model(&rs)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		query = query.Where("name ilike ?", "%"+params.Name+"%")
	}

	if params.GroupID != "" {
		query = query.Where("group_id = ?", params.GroupID)
	}

	if err := query.Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(rs))
	for i, r := range rs {
		result[i] = core.Model(r)
	}

	return result, nil
}

// Upsert upserts a role
func (s *RoleQueryStrategy) Upsert(m core.Model) error {
	_, err := s.DB.Model(m).
		OnConflict("(id) DO UPDATE").
		Set(`(
			"name",
			"group_id",
			"privilege"
		) = (
			?name,
			?group_id,
			?privilege)`).
		Insert()
	return err
}

// UpsertRelationship upserts a role relationship
func (s *RoleQueryStrategy) UpsertRelationship(e core.Entity, relation string) error {
	var query *orm.Query
	switch relation {
	default:
		return fmt.Errorf("Unknown relation for role: %s", relation)
	}

	_, err := query.Insert()
	return err
}
