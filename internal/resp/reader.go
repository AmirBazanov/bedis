package resp

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"

	l "bedis/pkg/logger"
)

var (
	ErrUnknownType = errors.New("unknown type of cmd")
	ErrInvalidSize = errors.New("invalid parse size")
	ErrInvalidCrlf = errors.New("invalid crlf")
)

const MAXSIZE = 512 * 1024 * 1024

type Reader struct {
	logger *slog.Logger
	reader *bufio.Reader
}

func NewReader(reader io.Reader, logger *slog.Logger) *Reader {
	logger = l.LoggerNotInitialized(logger)
	return &Reader{
		logger: logger,
		reader: bufio.NewReader(reader),
	}
}

func (r *Reader) Value() (*Value, error) {
	op := "reader.ReadValue"
	v, err := r.reader.ReadString('\n')
	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}

	switch Type(v[0]) {
	case BulkString:
		size, err := r.parseSize(v)
		if err != nil {
			return nil, err
		}
		r.logger.Info(op, slog.Int("reading data with size", size))
		r.logger.Info(op, slog.String("info", "RESP: BulkString"))
		return r.bulkString(size)
	case Array:
		size, err := r.parseSize(v)
		if err != nil {
			return nil, err
		}
		r.logger.Info(op, slog.Int("reading data with size", size))
		r.logger.Info(op, slog.String("info", "RESP: Array"))
		return r.array(size)
	case Integer:
		r.logger.Info(op, slog.String("info", "RESP: Integer"))
		return r.integer(v)
	case SimpleString:
		r.logger.Info(op, slog.String("info", "RESP: SimpleString"))
		return r.simpleString(v)
	case SimpleError:
		r.logger.Info(op, slog.String("info", "RESP: SimpleError"))
		return r.simpleError(v)
	default:
		return nil, ErrUnknownType
	}
}

func (r *Reader) bulkString(size int) (*Value, error) {
	op := "reader.BulkString"
	if size == -1 {
		r.logger.Warn(op, slog.String("info", " size of string is -1"))
		return &Value{Type: BulkString, Bytes: nil}, nil
	}
	if size < -1 {
		r.logger.Error(op, slog.String("error", "size less then -1"))
		return nil, ErrInvalidSize
	}
	buf := make([]byte, size)
	_, err := io.ReadFull(r.reader, buf)
	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}
	var crlf [2]byte
	if _, err := io.ReadFull(r.reader, crlf[:]); err != nil {
		r.logger.Error(op, slog.Any("error read crlf", err))
	}
	if crlf[0] != '\r' || crlf[1] != '\n' {
		r.logger.Error(op, slog.String("error", "ending not crlf"))
		return nil, ErrInvalidCrlf
	}
	return &Value{Type: BulkString, Bytes: buf}, nil
}

func (r *Reader) integer(data string) (*Value, error) {
	op := "reader.Integer"
	if len(data) < 3 {
		r.logger.Error(op, slog.Any("error", ErrInvalidSize))
		return nil, ErrInvalidSize
	}
	if data[len(data)-2:] != "\r\n" {
		r.logger.Error(op, slog.Any("error", ErrInvalidCrlf))
		return nil, ErrInvalidCrlf
	}
	var intVal int64
	intVal, err := strconv.ParseInt(data[1:len(data)-2], 10, 64)
	if err != nil {
		r.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}
	return &Value{Type: Integer, Integer: intVal}, nil
}

func (r *Reader) array(size int) (*Value, error) {
	op := "reader.Array"
	if size == -1 {
		r.logger.Error(op, slog.String("info", "arrays size is -1"))
		return &Value{Type: Array, IsNil: true}, nil
	}

	values := make([]*Value, 0, size)
	for range size {
		val, err := r.Value()
		if err != nil {
			r.logger.Error(op, slog.String("error", err.Error()))
			return nil, err
		}
		r.logger.Info(op, slog.String("info", "appending to array"))
		values = append(values, val)
	}
	return &Value{Type: Array, Array: values}, nil
}

func (r *Reader) simpleString(data string) (*Value, error) {
	op := "reader.SimpleString"
	if len(data) < 3 {
		r.logger.Error(op, slog.Any("error", ErrInvalidSize))
		return nil, ErrInvalidSize
	}
	if data[len(data)-2:] != "\r\n" {
		r.logger.Error(op, slog.Any("error", ErrInvalidCrlf))
		return nil, ErrInvalidCrlf
	}
	return &Value{Type: SimpleString, Bytes: []byte(data[1 : len(data)-2])}, nil
}

func (r *Reader) simpleError(data string) (*Value, error) {
	op := "reader.SimpleError"
	simpleStringValue, err := r.simpleString(data)
	if err != nil {
		r.logger.Error(op, slog.Any("error", err))
		return nil, err
	}
	return &Value{Type: SimpleError, Bytes: simpleStringValue.Bytes}, nil
}

// TODO: Add other types
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
	if size > MAXSIZE {

		r.logger.Error(op, slog.String("error", "size overflow"))
		return 0, ErrInvalidSize
	}
	if size < 0 {

		r.logger.Error(op, slog.String("error", "size is negative"))
		return 0, ErrInvalidSize
	}
	return size, nil
}
