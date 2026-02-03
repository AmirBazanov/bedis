package resp

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"strconv"

	l "bedis/pkg/logger"
	ws "bedis/pkg/writersticky"
)

var (
	ErrToWrite  = errors.New("unable to write in writer")
	ErrNilValue = errors.New("value is undefined")
	CRLF        = []byte{'\r', '\n'}
)

type Writer struct {
	writer *bufio.Writer
	logger *slog.Logger
	sw     *ws.WriterSticky
}

func NewWriter(writer io.Writer, logger *slog.Logger) *Writer {
	logger = l.LoggerNotInitialized(logger)
	bioWriter := bufio.NewWriter(writer)
	return &Writer{
		logger: logger,
		writer: bioWriter,
		sw:     &ws.WriterSticky{W: bioWriter, Err: nil},
	}
}

func (w *Writer) Value(value *Value) error {
	op := "writer.WriterValue"
	w.sw.Err = nil
	if value == nil {
		w.logger.Error(op, slog.Any("ERROR:", ErrNilValue))
		return ErrNilValue
	}
	switch value.Type {
	case SimpleString:
		return w.simpleString(value)
	case Integer:
		return w.integer(value)
	case SimpleError:
		return w.simpleError(value)
	case BulkString:
		return w.bulkString(value)
	case Array:
		return w.array(value)
	default:
		return ErrUnknownType
	}
}

func (w *Writer) simpleString(data *Value) error {
	op := "writer.SimpleString"
	w.sw.WriteByte(byte(data.Type))
	w.sw.Write(data.Bytes)
	w.sw.Write(CRLF)
	return w.handleErrOnWrite(w.sw.Err, "SimpleString", op)
}

func (w *Writer) integer(data *Value) error {
	op := "writer.Integer"
	w.sw.WriteByte(byte(data.Type))
	w.sw.WriteString(strconv.FormatInt(data.Integer, 10))
	w.sw.Write(CRLF)
	return w.handleErrOnWrite(w.sw.Err, "integer", op)
}

func (w *Writer) simpleError(data *Value) error {
	op := "writer.simpleError"

	w.sw.WriteByte(byte(data.Type))
	w.sw.Write(data.Bytes)
	w.sw.Write(CRLF)
	return w.handleErrOnWrite(w.sw.Err, "simpleError", op)
}

func (w *Writer) bulkString(data *Value) error {
	op := "writer.bulkString"

	w.sw.WriteByte(byte(data.Type))
	if data.Bytes == nil {
		w.sw.WriteString("-1")
		w.sw.Write(CRLF)
		return w.handleErrOnWrite(w.sw.Err, "bulkString", op)
	}
	w.sw.WriteString(strconv.Itoa(len(data.Bytes)))
	w.sw.Write(CRLF)
	w.sw.Write(data.Bytes)
	w.sw.Write(CRLF)
	return w.handleErrOnWrite(w.sw.Err, "bulkString", op)
}

func (w *Writer) array(data *Value) error {
	op := "writer.array"

	w.sw.WriteByte(byte(data.Type))
	if data.Array == nil {
		w.sw.WriteString("-1")
		w.sw.Write(CRLF)
		return w.handleErrOnWrite(w.sw.Err, "array", op)
	}
	w.sw.WriteString(strconv.Itoa(len(data.Array)))
	w.sw.Write(CRLF)
	if w.sw.Err != nil {
		w.logger.Error(op, slog.Any(ErrToWrite.Error(), w.sw.Err))
		return w.sw.Err
	}
	for i := range data.Array {
		err := w.Value(data.Array[i])
		if err != nil {
			w.logger.Error(op, slog.Any("in array err", err))
			return err
		}
	}
	return nil
}

func (w *Writer) handleErrOnWrite(err error, typ string, op string) error {
	if err != nil {

		w.logger.Error(op, slog.Any(ErrToWrite.Error(), err))
		return err
	}

	return nil
}

func (w *Writer) Flush() error {
	return w.writer.Flush()
}
