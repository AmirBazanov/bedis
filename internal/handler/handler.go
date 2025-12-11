package handler

import (
	"bedis/internal/storage"
	"fmt"
	"strings"
)

type Handler struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) Process(line string) (string, error) {
	parts := strings.Fields(line)
	command := strings.ToUpper(parts[0])
	fmt.Println("command:", command)
	switch command {
	case "SET":
		if len(parts) != 3 {
			return "", fmt.Errorf("wrong number of arguments: %s", line)
		}
		key := parts[1]
		value := []byte(parts[2])
		err := h.storage.Set(key, value)
		if err != nil {
			return "", err
		}
		return "OK", nil
	case "GET":
		if len(parts) != 2 {
			return "", fmt.Errorf("wrong number of arguments: %s", line)
		}
		key := parts[1]
		value, err := h.storage.Get(key)
		if err != nil || value == nil {
			return "", err
		}
		return string(value), nil
	default:
		return "", fmt.Errorf("unknown command: %s", command)

	}

}
