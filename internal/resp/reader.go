package resp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"log/slog"
	"strconv"
	"strings"
)

var (
	ErrUnknownType = errors.New("unknown type of cmd")
	ErrInvalidSize = errors.New("invalid parse size")
	ErrInvalidCrlf = errors.New("invalid crlf")
)

type Reader struct {
	logger *slog.Logger
	reader *bufio.Reader
}

func New(reader io.Reader, logger *slog.Logger) *Reader {
	op := "reader.New"
	if logger == nil {
		logger = slog.New(
			slog.NewTextHandler(io.Discard, nil),
		)
		log.Print(op + " no logger provided")
	}
	return &Reader{
		logger: logger,
		reader: bufio.NewReader(reader),
	}
}

func (r *Reader) ReadValue() (*Value, error) {
	op := "reader.ReadValue"
	v, err := r.reader.ReadString('\n')
	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}

	switch v[0] {
	case byte(BulkString):
		size, err := r.parseSize(v)
		if err != nil {
			return nil, err
		}
		r.logger.Info(op, slog.Int("reading data with size", size))
		r.logger.Info(op, slog.String("info", "RESP: BulkString"))
		return r.readBulk(size)
	case byte(Array):
		size, err := r.parseSize(v)
		if err != nil {
			return nil, err
		}
		r.logger.Info(op, slog.Int("reading data with size", size))
		r.logger.Info(op, slog.String("info", "RESP: Array"))
		return r.readArray(size)
	default:
		return nil, ErrUnknownType
	}
}

func (r *Reader) readBulk(size int) (*Value, error) {
	op := "reader.readBulk"
	if size == -1 {
		r.logger.Warn(op, slog.String("info", " size of string is -1"))
		return &Value{Type: BulkString}, nil
	}
	buf := make([]byte, size)
	_, err := io.ReadFull(r.reader, buf)
	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}
	crlf, err := r.reader.ReadString('\n')

	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}
	if crlf != "\r\n" {
		return nil, ErrInvalidCrlf
	}
	return &Value{Type: BulkString, Bytes: buf}, nil

}

func (r *Reader) readArray(size int) (*Value, error) {
	op := "reader.readArray"
	if size == -1 {
		r.logger.Error(op, slog.String("info", "arrays size is -1"))
		return &Value{Type: Array, IsNil: true}, nil
	}

	values := make([]*Value, 0, size)
	for i := 0; i < size; i++ {
		val, err := r.ReadValue()
		if err != nil {
			r.logger.Error(op, slog.String("error", err.Error()))
			return nil, err
		}
		r.logger.Info(op, slog.String("info", "appending to array"))
		values = append(values, val)
	}
	return &Value{Type: Array, Array: values}, nil
}

//TODO: Add other types

func (r *Reader) parseSize(sb string) (int, error) {
	op := "reader.parseSize"
	if len(sb) < 3 || !strings.HasSuffix(sb, "\r\n") {
		return 0, ErrInvalidSize
	}

	size, err := strconv.Atoi(string(sb[1 : len(sb)-2]))
	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return 0, ErrInvalidSize
	}
	return size, nil
}
