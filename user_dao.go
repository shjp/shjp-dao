package dao

import (
	"github.com/go-pg/pg"

	"github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

type userDAO struct {
	DB *pg.DB
}

type user struct {
	model.User

	tableName struct{} `sql:"select:users_full"`
}

// GetAll returns all users
func (o *userDAO) GetAll() ([]core.Model, error) {
	users := make([]*user, 0)

	if err := o.DB.Model(&users).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(users))
	for i, u := range users {
		result[i] = core.Model(u)
	}

	return result, nil
}

// GetOne returns one user
func (o *userDAO) GetOne(id string) (core.Model, error) {
	var u user
	var err error
	u.ID = id
	if err := o.DB.Model(&u).First(); err != nil {
		return nil, err
	}

	return &u, err
}

// Upsert upserts a user
func (o *userDAO) Upsert(m core.Model) error {
	return o.DB.Insert(m)
}
