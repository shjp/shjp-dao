package dao

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	core "github.com/shjp/shjp-core"
)

// ModelService is the http handler function factory for models
type ModelService struct {
	modelName string
	dao       modelDAO
}

// NewModelService returns a new model service
func NewModelService(modelName string, db *pg.DB) (*ModelService, error) {
	log.Printf("Initializing a new model service for %s...\n", modelName)
	var dao modelDAO
	switch modelName {
	case "group":
		dao = &groupDAO{DB: db}
	case "user":
		dao = &userDAO{DB: db}
	case "announcement":
		dao = &announcementDAO{DB: db}
	case "event":
		dao = &eventDAO{DB: db}
	default:
		return nil, fmt.Errorf("Model '%s' is not implemented", modelName)
	}

	return &ModelService{modelName: modelName, dao: dao}, nil
}

// HandleGetAll handles get all request
func (s *ModelService) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	models, err := s.dao.GetAll()
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error getting '%s' all models: %s", s.modelName, err))
		return
	}
	bytes, err := json.Marshal(models)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error serializing '%s' models: %s", s.modelName, err))
		return
	}
	fmt.Fprintf(w, string(bytes))
}

// HandleGetOne handles get one request
func (s *ModelService) HandleGetOne(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	model, err := s.dao.GetOne(id)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error getting '%s' one model: %s", s.modelName, err))
		return
	}
	bytes, err := json.Marshal(model)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error serializing '%s' one model: %s", s.modelName, err))
		return
	}
	fmt.Fprintf(w, string(bytes))
}

// HandleUpsert handles upsert request
func (s *ModelService) HandleUpsert(m core.Model) error {
	return s.dao.Upsert(m)
}
