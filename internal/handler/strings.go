package handler

import (
	"bedis/internal/resp"
	"bedis/internal/storage"
)

func handleSet(args []*resp.Value, s *storage.Storage) (*resp.Value, error) {
	key := string(args[0].Bytes)
	value := args[1].Bytes
	err := s.Set(key, value)
	if err != nil {
		return nil, err
	}
	return &resp.Value{Type: resp.SimpleString, Bytes: []byte("OK")}, nil
}

func handleGet(args []*resp.Value, s *storage.Storage) (*resp.Value, error) {
	key := string(args[0].Bytes)
	value, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	return &resp.Value{Type: resp.BulkString, Bytes: value}, nil
}
