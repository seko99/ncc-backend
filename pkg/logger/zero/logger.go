package zero

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

type Logger struct {
	logger zerolog.Logger
}

func (s *Logger) Info(msg string, fields ...interface{}) {
	s.logger.Info().Msgf(msg, fields...)
}

func (s *Logger) Warn(msg string, fields ...interface{}) {
	s.logger.Warn().Msgf(msg, fields...)
}

func (s *Logger) Error(msg string, fields ...interface{}) {
	s.logger.Error().Msgf(msg, fields...)
}

func (s *Logger) Debug(msg string, fields ...interface{}) {
	s.logger.Debug().Msgf(msg, fields...)
}

func (s *Logger) Trace(msg string, fields ...interface{}) {
	s.logger.Trace().Msgf(msg, fields...)
}

func NewLogger(l ...zerolog.Level) *Logger {
	level := zerolog.InfoLevel

	if len(l) > 0 {
		level = l[0]
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    true,
		TimeFormat: "2006/01/02 - 15:04:05.000000",
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
		FormatLevel: func(i interface{}) string {
			if l, ok := i.(string); ok {
				switch l {
				case zerolog.LevelTraceValue:
					return "[TRACE]"
				case zerolog.LevelDebugValue:
					return "[DEBUG]"
				case zerolog.LevelInfoValue:
					return "[INFO ]"
				case zerolog.LevelWarnValue:
					return "[WARN ]"
				case zerolog.LevelErrorValue:
					return "[ERROR]"
				case zerolog.LevelFatalValue:
					return "[FATAL]"
				case zerolog.LevelPanicValue:
					return "[PANIC]"
				default:
					return "[UNKNW]"
				}
			} else {
				return "[UNKNW]"
			}
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s", i)
		},
	}

	zerolog.SetGlobalLevel(level)
	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &Logger{
		logger: logger,
	}
}
