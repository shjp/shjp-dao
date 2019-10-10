package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"

	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// GroupQueryStrategy implements QueryStrategy for groups
type GroupQueryStrategy struct {
	*pg.DB
}

type group struct {
	model.Group

	tableName struct{} `sql:"select:groups_full"`
}

type groupMembership struct {
	core.GroupMembership

	tableName struct{} `sql:"groups_users"`
}

// ModelName outputs this model's name
func (s *GroupQueryStrategy) ModelName() string {
	return "group"
}

// GetAll returns all groups
func (s *GroupQueryStrategy) GetAll() ([]core.Model, error) {
	groups := make([]*group, 0)

	err := s.DB.Model(&groups).Column("*").Select()
	if err != nil {
		return nil, err
	}

	result := make([]core.Model, len(groups))
	for i, g := range groups {
		result[i] = core.Model(g)
	}

	return result, nil
}

// GetOne returns one group
func (s *GroupQueryStrategy) GetOne(id string) (core.Model, error) {
	var g group
	g.ID = id
	if err := s.DB.Model(&g).WherePK().First(); err != nil {
		return nil, err
	}

	return &g, nil
}

// Search finds all group meeting the criteria given by the payload
func (s *GroupQueryStrategy) Search(payload []byte) ([]core.Model, error) {
	var params group
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	gs := make([]*group, 0)

	query := s.DB.Model(&gs)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if params.Name != "" {
		query = query.Where("name ilike ?", "%"+params.Name+"%")
	}

	if err := query.Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(gs))
	for i, g := range gs {
		result[i] = core.Model(g)
	}

	return result, nil
}

// Upsert upserts a group
func (s *GroupQueryStrategy) Upsert(m core.Model) error {
	_, err := s.DB.Model(m).
		OnConflict("(id) DO UPDATE").
		Set(`(
			"name",
			"description",
			"image_url"
		) = (
			?name,
			?description,
			?image_url)`).
		Insert()
	return err
}

// UpsertRelationship upserts a group relationship
func (s *GroupQueryStrategy) UpsertRelationship(e core.Entity, relation string) error {
	var dbModel core.Entity
	var query *orm.Query
	switch relation {
	case "group_membership":
		dbModel = &groupMembership{
			GroupMembership: *e.(*core.GroupMembership),
		}
		query = s.DB.Model(dbModel).
			OnConflict("(group_id, user_id) DO UPDATE").
			Set(`(
				"group_id",
				"user_id",
				"role_id",
				"status"
			) = (
				?group_id,
				?user_id,
				?role_id,
				?status
			)`)
	default:
		return fmt.Errorf("Unknown relation for group: %s", relation)
	}

	_, err := query.Insert()
	return err
}
