package log

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/configuration"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

const logFile = "./log/ara-server.log"

var logger zerolog.Logger

func Info(ctx context.Context, metadata interface{}, error error, message string) {
	var ctxID string
	if ctx != nil {
		if id, ok := ctx.Value(constants.CtxKeyCtxID).(string); ok {
			ctxID = id
		}
	}
	logger.Info().Timestamp().Str("ctx_id", ctxID).Interface("metadata", metadata).Err(error).Msg(message)
}

func Error(ctx context.Context, metadata interface{}, error error, message string) {
	var ctxID string
	if ctx != nil {
		if id, ok := ctx.Value(constants.CtxKeyCtxID).(string); ok {
			ctxID = id
		}
	}
	logger.Error().Timestamp().Str("ctx_id", ctxID).Interface("metadata", metadata).Err(error).Msg(message)
}

func Fatal(ctx context.Context, metadata interface{}, error error, message string) {
	var ctxID string
	if ctx != nil {
		if id, ok := ctx.Value(constants.CtxKeyCtxID).(string); ok {
			ctxID = id
		}
	}
	logger.Fatal().Timestamp().Str("ctx_id", ctxID).Interface("metadata", metadata).Err(error).Msg(message)
}

func NewLogger(config configuration.Config) {
	err := os.MkdirAll(filepath.Dir(logFile), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	var writer io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: "2006-01-02 15:04:05",
	}

	if !config.IsDevelopment() {
		writer = file
	}

	logger = zerolog.New(writer).With().Logger()
}
