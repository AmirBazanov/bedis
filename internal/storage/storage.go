/*
Package storage
*/
package storage

import (
	"errors"
	"sync"
)

type Storage struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

func New() *Storage {
	return &Storage{
		data: make(map[string][]byte),
	}
}

func (s *Storage) Set(key string, value []byte) error {
	s.mutex.Lock()
	data := make([]byte, len(value))
	copy(data, value)
	s.data[key] = data
	s.mutex.Unlock()
	return nil
}

func (s *Storage) Get(key string) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ok := s.data[key]
	if !ok {
		return nil, errors.New("NO SUCH DATA")
	}
	data := make([]byte, len(item))
	copy(data, item)
	return data, nil
}
