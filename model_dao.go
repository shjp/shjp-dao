package dao

import core "github.com/shjp/shjp-core"

type modelDAO interface {
	GetAll() ([]core.Model, error)
	GetOne(string) (core.Model, error)
	Search([]byte) ([]core.Model, error)
	Upsert(core.Model) error
}
