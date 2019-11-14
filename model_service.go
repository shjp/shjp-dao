package dao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// ModelService is the http handler function factory for models
type ModelService struct {
	QueryStrategy
}

// NewModelService returns a new model service
func NewModelService(queryStrategy QueryStrategy) *ModelService {
	log.Printf("Initializing a new model service for %s...\n", queryStrategy.ModelName())
	return &ModelService{QueryStrategy: queryStrategy}
}

// HandleGetAll handles get all request
func (s *ModelService) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	ms, err := s.QueryStrategy.GetAll()
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error getting '%s' all models: %s", s.ModelName(), err))
		return
	}
	bytes, err := json.Marshal(ms)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error serializing '%s' models: %s", s.ModelName(), err))
		return
	}
	fmt.Fprintf(w, string(bytes))
}

// HandleGetOne handles get one request
func (s *ModelService) HandleGetOne(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	m, err := s.QueryStrategy.GetOne(id)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error getting '%s' one model: %s", s.ModelName(), err))
		return
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error serializing '%s' one model: %s", s.ModelName(), err))
		return
	}
	fmt.Fprintf(w, string(bytes))
}

// HandleSearch handles the request to get all models meeting the criteria
func (s *ModelService) HandleSearch(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error reading body for HandleFind"))
		return
	}
	ms, err := s.QueryStrategy.Search(payload)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error finding '%s' models: %s", s.ModelName(), err))
		return
	}
	bytes, err := json.Marshal(ms)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error serializing '%s' models: %s", s.ModelName(), err))
		return
	}
	fmt.Fprintf(w, string(bytes))
}

// HandleCreate handles the REST API request to create the model object
func (s *ModelService) HandleCreate(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error reading body for HandleRestUpsert"))
		return
	}
	m, err := s.read(payload)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error deserializing '%s' model payload: %s", s.ModelName(), err))
		return
	}
	if err = s.HandleUpsert(m); err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error upserting '%s' model: %s", s.ModelName(), err))
		return
	}
}

// HandleUpsert handles upsert request
func (s *ModelService) HandleUpsert(m core.Model) error {
	return s.QueryStrategy.Upsert(m)
}

// HandleUpsertRelationship handles upsert request for relationship
func (s *ModelService) HandleUpsertRelationship(e core.Entity, relation string) error {
	return s.QueryStrategy.UpsertRelationship(e, relation)
}

func (s *ModelService) read(b []byte) (core.Model, error) {
	var m core.Model
	switch s.ModelName() {
	case "announcement":
		m = &model.Announcement{}
	case "event":
		m = &model.Event{}
	case "group":
		m = &model.Group{}
	case "role":
		m = &model.Role{}
	case "user":
		m = &model.User{}
	}
	if err := json.Unmarshal(b, m); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling model")
	}
	return m, nil
}
