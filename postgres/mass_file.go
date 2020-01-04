package postgres

import (
	"github.com/go-pg/pg"

	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// MassFileQueryStrategy implements QueryStrategy for mass files
type MassFileQueryStrategy struct {
	*pg.DB
}

type massFile struct {
	model.MassFile
}

// ModelName outputs this model's name
func (s *MassFileQueryStrategy) ModelName() string {
	return "mass_file"
}

// GetAll returns all roles
func (s *MassFileQueryStrategy) GetAll() ([]core.Model, error) {
	massFiles := make([]*massFile, 0)

	err := s.DB.Model(&massFiles).Select()
	if err != nil {
		return nil, err
	}

	result := make([]core.Model, len(massFiles))
	for i, f := range massFiles {
		result[i] = core.Model(f)
	}

	return result, nil
}

// GetOne returns one role
func (s *MassFileQueryStrategy) GetOne(id string) (core.Model, error) {
	return nil, errNotImplemented
}

// Search finds all roles meeting the criteria given by the payload
func (s *MassFileQueryStrategy) Search(payload []byte) ([]core.Model, error) {
	return nil, errNotImplemented
}

// Upsert upserts a role
func (s *MassFileQueryStrategy) Upsert(m core.Model) error {
	return errNotImplemented
}

// UpsertRelationship upserts a role relationship
func (s *MassFileQueryStrategy) UpsertRelationship(e core.Entity, relation string) error {
	return errNotImplemented
}
