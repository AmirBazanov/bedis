/*
Package storage
*/
package storage

import (
	"errors"
	"io"
	"log/slog"
	"sync"
)

// TODO: Errors
type Storage struct {
	data   map[string][]byte
	mutex  sync.RWMutex
	logger *slog.Logger
}

func New(logger *slog.Logger) *Storage {
	if logger == nil {
		logger = slog.New(
			slog.NewTextHandler(io.Discard, nil),
		)
	}

	return &Storage{
		data:   make(map[string][]byte),
		logger: logger,
	}
}

func (s *Storage) Set(key string, value []byte) error {
	op := "storage.Set"
	s.mutex.Lock()
	data := make([]byte, len(value))
	s.logger.Info(op, "Setting data", slog.String("key", key))
	copy(data, value)
	s.data[key] = data
	s.mutex.Unlock()
	return nil
}

func (s *Storage) Get(key string) ([]byte, error) {
	op := "storage.Get"
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ok := s.data[key]
	s.logger.Info(op, "Getting data", slog.String("key", key))
	if !ok {
		return nil, errors.New("NO SUCH DATA")
	}
	data := make([]byte, len(item))
	copy(data, item)
	return data, nil
}
