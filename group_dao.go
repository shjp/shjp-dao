package dao

import (
	"github.com/go-pg/pg"

	"github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

type groupDAO struct {
	DB *pg.DB
}

type group struct {
	model.Group

	tableName struct{} `sql:"groups"`
}

// GetAll returns all groups
func (o *groupDAO) GetAll() ([]core.Model, error) {
	groups := make([]*group, 0)

	err := o.DB.Model(&groups).
		ColumnExpr(`"group".*`).
		ColumnExpr("COALESCE(json_agg(users) FILTER (WHERE users.id IS NOT NULL), '[]') AS members").
		Join(`LEFT JOIN groups_users AS gu ON gu.group_id = "group".id`).
		Join("LEFT JOIN users ON users.id = gu.user_id").
		Group("group.id").
		Select()
	if err != nil {
		return nil, err
	}

	result := make([]core.Model, len(groups))
	for i, g := range groups {
		result[i] = core.Model(g)
	}

	return result, nil
}

// Upsert upserts a group
func (o *groupDAO) Upsert(m core.Model) error {
	_, err := o.DB.Model(m).
		OnConflict("(id) DO UPDATE").
		Set(`(
			name,
			description,
			image_url
		) = (
			?name,
			?description,
			?image_url)`).
		Insert()
	return err
}
