package interfaces

import "code.evixo.ru/ncc/ncc-backend/pkg/events"

//go:generate mockgen -destination=mocks/mock_events.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/events/interfaces Events

type Events interface {
	Run()
	PublishProtoEvent(event events.Event) error
	PublishEvent(event events.Event) error
	PublishRequest(event events.Event, responseType string, cb func(event events.Event, params ...interface{}), params ...interface{}) error
	SubscribeOnProtoEvent(eventType string, cb func(event events.Event)) error
	SubscribeOnEvent(eventType string, cb func(event events.Event)) error
	SubscribeOnBroadcast(eventType string, cb func(event events.Event)) error
}
