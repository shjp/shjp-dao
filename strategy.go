package dao

import core "github.com/shjp/shjp-core"

// QueryStrategy provides interface for DB query operations
type QueryStrategy interface {
	// Outputs the model name
	ModelName() string

	// Query operations
	GetAll() ([]core.Model, error)
	GetOne(string) (core.Model, error)
	Search([]byte) ([]core.Model, error)
	Upsert(core.Model) error
	UpsertRelationship(core.Entity, string) error
}
