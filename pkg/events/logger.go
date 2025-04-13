package events

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"github.com/ThreeDotsLabs/watermill"
)

type EventLogger struct {
	log logger.Logger
}

//todo: implement log levels

func (l EventLogger) Error(msg string, err error, fields watermill.LogFields) {
	l.log.Error(msg)
}

func (l EventLogger) Info(msg string, fields watermill.LogFields) {
	l.log.Debug(msg)
}

func (l EventLogger) Debug(msg string, fields watermill.LogFields) {
	//todo: loglevel
	//l.log.Debug(msg)
}

func (l EventLogger) Trace(msg string, fields watermill.LogFields) {
	l.log.Trace(msg)
}

func (l EventLogger) With(fields watermill.LogFields) watermill.LoggerAdapter { return l }

func NewEventLogger(log logger.Logger) watermill.LoggerAdapter {
	return &EventLogger{
		log: log,
	}
}
