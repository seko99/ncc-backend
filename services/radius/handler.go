package radius

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events/interfaces"
	"log"
)

type Handler struct {
	*log.Logger
	Verbose bool
	events  interfaces.Events
}

func NewHanlder(events interfaces.Events) *Handler {
	return &Handler{
		events: events,
	}
}
