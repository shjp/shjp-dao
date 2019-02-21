package dao

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	core "github.com/shjp/shjp-core"
	"github.com/shjp/shjp-queue"
)

const (
	origin = "dao"
)

// AsyncService is responsible for handling asynchronous requests
type AsyncService struct {
	consumer *queue.Consumer
	producer *queue.Producer
	exchange string

	services map[string]*ModelService
}

// NewAsyncService instantiates a new AsyncService
func NewAsyncService(queueHost, queueUser, queueExchange string, svcs ...*ModelService) (*AsyncService, error) {
	log.Printf("Initializing a new AsyncService... | Host: %s | User: %s | Exchange: %s\n", queueHost, queueUser, queueExchange)
	consumer, err := queue.NewConsumer(queueHost, queueUser, origin, core.IntentRequest)
	if err != nil {
		return nil, err
	}
	producer, err := queue.NewProducer(queueHost, queueUser)
	if err != nil {
		return nil, err
	}

	services := make(map[string]*ModelService)
	for _, service := range svcs {
		services[service.ModelName()] = service
	}

	return &AsyncService{
		consumer,
		producer,
		queueExchange,
		services,
	}, nil
}

// Listen starts consuming messages from the queue and executes the handler function upon receipt
func (s *AsyncService) Listen() error {
	return s.consumer.Consume(true, s.handle)
}

func (s *AsyncService) handle(msg *core.Message) {
	if msg.Intent != core.IntentRequest {
		log.Println("Ignoring non-request intent message")
		return
	}
	var messageHandler func(*core.Message) error
	switch msg.Type {
	case core.ModelType:
		messageHandler = s.handleModelMessage
	case core.RelationshipType:
		messageHandler = s.handleRelationshipMessage
	default:
		log.Println("Ignoring messages with other than model and relationship types")
		return
	}

	if err := messageHandler(msg); err != nil {
		log.Println(err)
		return
	}

	s.postProcess()
}

func (s *AsyncService) handleModelMessage(msg *core.Message) error {
	ms, ok := s.services[msg.Subtype]
	if !ok {
		// TODO: Send error message to the queue
		//s.producer.Publish()
		return fmt.Errorf("Cannot handle model message with unknown subtype: %s", msg.Subtype)
	}

	entity, err := msg.ExtractEntity()
	if err != nil {
		// TODO: send error message
		return errors.Wrap(err, "Cannot extract entity from message")
	}

	m, ok := entity.(core.Model)
	if !ok {
		// TODO: send error message
		return errors.Wrap(err, "Cannot convert the entity to model")
	}

	switch msg.OperationType {
	case core.UpsertOperation:
		err = ms.HandleUpsert(m)
	default:
		err = fmt.Errorf("Unrecognized operation type %s", msg.OperationType)
	}

	if err != nil {
		// TODO: send error message
		return errors.Wrap(err, "Operation not successful")
	}

	return nil
}

func (s *AsyncService) handleRelationshipMessage(msg *core.Message) error {
	var serviceName string
	switch msg.Subtype {
	case "group_membership":
		serviceName = "group"
	default:
		return fmt.Errorf("Cannot handle relationship message with unknown subtype: %s", msg.Subtype)
	}

	ms, ok := s.services[serviceName]
	if !ok {
		// TODO: Send error message to the queue
		return fmt.Errorf("Cannot handle relationship message with unknown service name: %s", serviceName)
	}

	entity, err := msg.ExtractEntity()
	if err != nil {
		// TODO: send error message
		return errors.Wrap(err, "Cannot extract entity from message")
	}

	switch msg.OperationType {
	case core.UpsertOperation:
		err = ms.HandleUpsertRelationship(entity, msg.Subtype)
	default:
		err = fmt.Errorf("Unrecognized operation type %s", msg.OperationType)
	}

	if err != nil {
		// TODO: send error message
		return errors.Wrap(err, "Operation not successful")
	}

	return nil
}

func (s *AsyncService) postProcess() {

}
