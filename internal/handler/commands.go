package handler

import (
	"bedis/internal/resp"
	"bedis/internal/storage"
)

type Command struct {
	Handler func(args []*resp.Value, s *storage.Storage) (*resp.Value, error)
	MinArgs int
	Name    string
}
