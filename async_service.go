package dao

import (
	"fmt"
	"log"

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
		services[service.modelName] = service
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
	if msg.Type != core.ModelType {
		log.Println("Ignoring non-model type message")
		return
	}
	ms, ok := s.services[msg.Subtype]
	if !ok {
		log.Println("Cannot handle model message with unknown subtype: ", msg.Subtype)
		// TODO: Send error message to the queue
		//s.producer.Publish()
		return
	}

	entity, err := msg.ExtractEntity()
	if err != nil {
		log.Println("Cannot extract entity from message:", err)
		// TODO: send error message
		return
	}

	m, ok := entity.(core.Model)
	if !ok {
		log.Println("Cannot convert the entity to model")
		// TODO: send error message
		return
	}

	switch msg.OperationType {
	case core.UpsertOperation:
		err = ms.HandleUpsert(m)
	default:
		err = fmt.Errorf("Unrecognized operation type %s", msg.OperationType)
	}

	if err != nil {
		log.Println("Operation not successful:", err)
		// TODO: send error message
		return
	}

	s.postProcess()
}

func (s *AsyncService) postProcess() {

}

/*func (s *AsyncService) chooseModelService(modelType string) (*Model, error) {
	switch modelType {
	case "announcement":
		return s.announcementService, nil
	case "event":
		return s.eventService, nil
	case "group":
		return s.groupService, nil
	case "user":
		return s.userService, nil
	default:
		return nil, fmt.Errorf("Unrecognized model service type: %s", modelType)
	}
}*/
