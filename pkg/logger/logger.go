package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

// Добавляем аргумент w io.Writer
func New(level string, w io.Writer) *Logger {
	var l zerolog.Level
	switch level {
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	// Если w не передан, используем стандартный вывод
	if w == nil {
		w = os.Stdout
	}

	skipFrameCount := 1
	return &Logger{
		// Вместо os.Stdout теперь используем w
		logger: zerolog.New(w).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger(),
	}
}

func (l *Logger) Info(message string, args ...interface{}) { l.logger.Info().Msgf(message, args...) }
func (l *Logger) Error(err error, message string, args ...interface{}) {
	l.logger.Error().Err(err).Msgf(message, args...)
}
func (l *Logger) Debug(message string, args ...interface{}) { l.logger.Debug().Msgf(message, args...) }
