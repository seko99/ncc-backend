package events

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	DefaultTopic = "events"

	BroadcastEvents = "broadcast.events"

	MaxReconnects    = 3
	ReconnectTimeout = 1 * time.Second
)

type Publisher struct {
	Id string `json:"id"`
}

type Event struct {
	Ts        time.Time              `json:"ts"`
	Publisher Publisher              `json:"publisher"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Message   proto.Message
}

type EventHandler struct {
	Id   string
	Type string
	cb   func(event Event)
}

type Events struct {
	cfg        *config.Config
	log        logger.Logger
	id         string
	config     amqp.Config
	publisher  *amqp.Publisher
	subscriber *amqp.Subscriber
	events     []*EventHandler
	ctx        context.Context
	Topic      string
	wlog       watermill.LoggerAdapter
	metrics    struct {
		sync.RWMutex
		events uint
	}
}

type RequestCallback func(event Event, params ...interface{})

func (s *Events) Run() {
	go func() {
		for range time.Tick(time.Second) {
			s.metrics.Lock()
			s.log.Trace("events=%d", s.metrics.events)
			s.metrics.events = 0
			s.metrics.Unlock()
		}
	}()
}

func (s *Events) PublishProtoEvent(event Event) error {
	b, _ := json.Marshal(Event{
		Ts: time.Now(),
		Publisher: Publisher{
			Id: s.id,
		},
		Type:    event.Type,
		Payload: event.Payload,
	})

	msg := message.NewMessage(watermill.NewUUID(), b)
	err := s.publisher.Publish(s.Topic, msg)
	if err != nil {
		s.log.Error("Publish error: %v", err)
	}
	return nil
}

func NewEvent(eventType string, payload interface{}) Event {
	var p map[string]interface{}

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return Event{}
		}
		err = json.Unmarshal(b, &p)
		if err != nil {
			return Event{}
		}
	}

	return Event{
		Type:    eventType,
		Payload: p,
	}
}

func (s *Events) PublishEvent(event Event) error {
	b, _ := json.Marshal(Event{
		Ts: time.Now(),
		Publisher: Publisher{
			Id: s.id,
		},
		Type:    event.Type,
		Payload: event.Payload,
	})

	msg := message.NewMessage(watermill.NewUUID(), b)
	err := s.publisher.Publish(s.Topic, msg)
	if err != nil {
		s.log.Error("Publish error: %v", err)
	}
	return nil
}

func (s *Events) PublishRequest(event Event, responseType string, cb func(event Event, params ...interface{}), params ...interface{}) error {
	b, _ := json.Marshal(Event{
		Ts: time.Now(),
		Publisher: Publisher{
			Id: s.id,
		},
		Type:    event.Type,
		Payload: event.Payload,
	})

	topic := s.Topic + "-" + responseType + "-" + s.id

	config := s.config
	config.Queue.AutoDelete = true
	subscriber, err := s.newSubscriber(config)
	if err != nil {
		return err
	}

	messages, err := subscriber.Subscribe(s.ctx, topic)
	if err != nil {
		panic(err)
	}

	go func(<-chan *message.Message) {
		for msg := range messages {
			var e Event
			err = json.Unmarshal(msg.Payload, &e)
			if err != nil {
				s.log.Error("Can't unmarshal event: %+v", err)
				msg.Nack()
			}

			s.log.Trace("Event: %+v", e)

			cb(e, params)
			msg.Ack()
			err = subscriber.Close()
			if err != nil {
				s.log.Error("Can't close subscriber: %v", err)
			}
		}
	}(messages)

	outMsg := message.NewMessage(watermill.NewUUID(), b)
	topic = s.Topic + "-" + event.Type
	err = s.publisher.Publish(topic, outMsg)
	if err != nil {
		s.log.Error("Publish error: %v", err)
	}

	return nil
}

func (s *Events) SubscribeOnProtoEvent(eventType string, cb func(event Event)) error {
	s.events = append(s.events, &EventHandler{
		Id:   s.id,
		Type: eventType,
		cb:   cb,
	})

	if len(s.events) > 0 && s.subscriber == nil {
		config := s.config
		config.Queue.AutoDelete = true
		subscriber, err := s.newSubscriber(config)
		if err != nil {
			return err
		}

		s.subscriber = subscriber

		messages, err := s.subscriber.Subscribe(s.ctx, s.Topic)
		if err != nil {
			panic(err)
		}
		go func(<-chan *message.Message) {
			for msg := range messages {
				var e Event

				err := json.Unmarshal(msg.Payload, &e)
				if err != nil {
					s.log.Error("Can't unmarshal event: %+v", err)
					msg.Nack()
				}

				s.log.Trace("Event: %+v, events: %+v", e, s.events)

				for _, registeredEvent := range s.events {
					if e.Type == registeredEvent.Type {
						re := registeredEvent
						go re.cb(e)
						msg.Ack()
					}
				}
			}

		}(messages)
	}

	return nil
}

func (s *Events) SubscribeOnEvent(eventType string, cb func(event Event)) error {
	s.events = append(s.events, &EventHandler{
		Id:   s.id,
		Type: eventType,
		cb:   cb,
	})

	if len(s.events) > 0 && s.subscriber == nil {
		c := s.config
		c.Queue.AutoDelete = true
		subscriber, err := s.newSubscriber(c)
		if err != nil {
			return err
		}

		s.subscriber = subscriber

		messages, err := s.subscriber.Subscribe(s.ctx, s.Topic)
		if err != nil {
			panic(err)
		}
		go func(<-chan *message.Message) {
			for msg := range messages {
				s.metrics.Lock()
				s.metrics.events++
				s.metrics.Unlock()

				var e Event

				err := json.Unmarshal(msg.Payload, &e)
				if err != nil {
					s.log.Error("Can't unmarshal event: %+v", err)
					msg.Nack()
				}

				s.log.Trace("Event: %+v, events: %+v", e, s.events)

				processed := false
				for _, registeredEvent := range s.events {
					if e.Type == registeredEvent.Type {
						re := registeredEvent
						go re.cb(e)
						processed = true
					}
				}
				if processed {
					msg.Ack()
				}
			}

		}(messages)
	}

	return nil
}

func (s *Events) SubscribeOnBroadcast(eventType string, cb func(event Event)) error {
	s.events = append(s.events, &EventHandler{
		Id:   s.id,
		Type: eventType,
		cb:   cb,
	})

	c := s.config
	c.Exchange = amqp.ExchangeConfig{
		GenerateName: func(topic string) string {
			return fmt.Sprintf("%s", topic)
		},
		Type:        "fanout",
		Durable:     true,
		AutoDeleted: true,
	}
	c.Queue.GenerateName = func(topic string) string {
		return fmt.Sprintf("%s-%s", topic, uuid.NewString())
	}
	c.Queue.AutoDelete = true
	subscriber, err := s.newSubscriber(c)
	if err != nil {
		return err
	}

	messages, err := subscriber.Subscribe(s.ctx, s.Topic)
	if err != nil {
		panic(err)
	}
	go func(<-chan *message.Message) {
		for msg := range messages {
			s.metrics.Lock()
			s.metrics.events++
			s.metrics.Unlock()

			var e Event

			err := json.Unmarshal(msg.Payload, &e)
			if err != nil {
				s.log.Error("Can't unmarshal event: %+v", err)
				msg.Nack()
			}

			s.log.Trace("Event: %+v, events: %+v", e, s.events)

			if e.Type == eventType {
				go cb(e)
			}
			msg.Ack()
		}

	}(messages)

	return nil
}

func (s *Events) SubscribeOnRequest(eventType, responseType string, cb func(event Event) map[string]interface{}) error {

	c := s.config
	c.Queue.AutoDelete = true
	subscriber, err := s.newSubscriber(c)
	if err != nil {
		return err
	}

	messages, err := subscriber.Subscribe(s.ctx, s.Topic+"-"+eventType)
	if err != nil {
		panic(err)
	}
	go func(<-chan *message.Message) {
		for msg := range messages {
			var e Event
			err := json.Unmarshal(msg.Payload, &e)
			if err != nil {
				s.log.Error("Can't unmarshal event: %+v", err)
				msg.Nack()
			}

			s.log.Trace("Event: %+v, events: %+v", e, s.events)

			if e.Type == eventType {
				payload := cb(e)
				b, _ := json.Marshal(Event{
					Publisher: Publisher{
						Id: s.id,
					},
					Type:    responseType,
					Payload: payload,
				})

				topic := s.Topic + "-" + responseType + "-" + e.Publisher.Id

				msg := message.NewMessage(watermill.NewUUID(), b)
				err = s.publisher.Publish(topic, msg)
				if err != nil {
					s.log.Error("Publish error: %v", err)
				}
			}
			msg.Ack()
		}

	}(messages)

	return nil
}

func (s *Events) newPublisher(config amqp.Config) (*amqp.Publisher, error) {
	var publisher *amqp.Publisher
	var err error

	for i := 0; i < MaxReconnects; i++ {
		publisher, err = amqp.NewPublisher(config, s.wlog)
		if err == nil {
			break
		}

		s.log.Error("Can't connect to AMQP, reconnecting: %v", err)
		time.Sleep(ReconnectTimeout)
	}

	return publisher, err
}

func (s *Events) newSubscriber(config amqp.Config) (*amqp.Subscriber, error) {
	var subscriber *amqp.Subscriber
	var err error

	for i := 0; i < MaxReconnects; i++ {
		subscriber, err = amqp.NewSubscriber(
			config,
			s.wlog,
		)
		if err == nil {
			break
		}

		s.log.Error("Can't create subscriber: %v", err)
		time.Sleep(ReconnectTimeout)
	}
	return subscriber, err
}

func NewEvents(cfg *config.Config, log logger.Logger, id string, topics ...string) (*Events, error) {
	amqpConfig := amqp.NewDurableQueueConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Queue.User,
		cfg.Queue.Password,
		cfg.Queue.Host,
		cfg.Queue.Port))

	wlog := NewEventLogger(log)

	amqpConfig.Publish.ChannelPoolSize = 16

	r := &Events{
		cfg:    cfg,
		log:    log,
		id:     id,
		config: amqpConfig,
		ctx:    context.Background(),
		wlog:   wlog,
	}

	var err error

	r.publisher, err = r.newPublisher(amqpConfig)
	if err != nil {
		return nil, err
	}

	if len(topics) > 0 {
		r.Topic = topics[0]
	} else {
		r.Topic = DefaultTopic
	}

	return r, nil
}

func NewBroadcastEvents(cfg *config.Config, log logger.Logger, id string, exchange string) (*Events, error) {
	amqpConfig := amqp.NewDurableQueueConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Queue.User,
		cfg.Queue.Password,
		cfg.Queue.Host,
		cfg.Queue.Port))

	wlog := NewEventLogger(log)

	amqpConfig.Publish.ChannelPoolSize = 16
	amqpConfig.Exchange = amqp.ExchangeConfig{
		GenerateName: func(topic string) string {
			return exchange
		},
		Type:        "fanout",
		Durable:     true,
		AutoDeleted: true,
	}

	r := &Events{
		cfg:    cfg,
		log:    log,
		id:     id,
		config: amqpConfig,
		ctx:    context.Background(),
		wlog:   wlog,
		Topic:  exchange,
	}

	var err error
	r.publisher, err = r.newPublisher(amqpConfig)
	if err != nil {
		return nil, err
	}

	return r, nil
}
