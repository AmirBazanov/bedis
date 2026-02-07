package handler

import (
	"errors"
	"log/slog"
	"strings"

	"bedis/internal/resp"
	"bedis/internal/storage"
	l "bedis/pkg/logger"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
	ErrWrongArgs      = errors.New("wrong number of arguments")
	ErrKeyNotFound    = errors.New("key not found")
)

type Handler struct {
	storage *storage.Storage
	logger  *slog.Logger
}

var CommandRegistry = map[string]Command{
	"SET": {
		Name:    "SET",
		MinArgs: 2,
		Handler: handleSet,
	},
	"GET": {
		Name:    "GET",
		MinArgs: 1,
		Handler: handleGet,
	},
}

func New(storage *storage.Storage, logger *slog.Logger) *Handler {
	logger = l.LoggerNotInitialized(logger)
	return &Handler{
		storage: storage,
		logger:  logger,
	}
}

func (h *Handler) Process(cmd *resp.Value) (*resp.Value, error) {
	op := "handler.Process"
	if cmd.Type != resp.Array || len(cmd.Array) == 0 {
		h.logger.Error(op, slog.Any("error", ErrUnknownCommand))
		return nil, ErrUnknownCommand
	}

	cmdName := strings.ToUpper(string(cmd.Array[0].Bytes))
	args := cmd.Array[1:]

	command, ok := CommandRegistry[cmdName]
	if !ok {
		return nil, ErrUnknownCommand
	}

	if len(args) < command.MinArgs {
		return nil, ErrWrongArgs
	}

	return command.Handler(args, h.storage)
}

func errorValue(msg string) *resp.Value {
	return &resp.Value{
		Type:  resp.SimpleError,
		Bytes: []byte(msg),
	}
}
