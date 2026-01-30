package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strconv"

	l "bedis/pkg/logger"
)

var (
	ErrToWrite  = errors.New("unable to write in writer")
	ErrNilValue = errors.New("value is undefined")
	CRLF        = []byte{'\r', '\n'}
)

type Writer struct {
	writer *bufio.Writer
	logger *slog.Logger
}

func NewWriter(writer io.Writer, logger *slog.Logger) *Writer {
	logger = l.LoggerNotInitialized(logger)
	return &Writer{
		logger: logger,
		writer: bufio.NewWriter(writer),
	}
}

func (w *Writer) Value(value *Value) error {
	op := "writer.WriterValue"
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
	default:
		return ErrUnknownType
	}
}

func (w *Writer) simpleString(data *Value) error {
	op := "writer.SimpleString"
	var buf bytes.Buffer
	buf.WriteByte(byte(data.Type))
	buf.Write(data.Bytes)
	buf.Write(CRLF)
	b, err := w.writer.Write(buf.Bytes())
	return w.handleErrOnWrite(err, b, "simpleString", op)
}

func (w *Writer) integer(data *Value) error {
	op := "writer.Integer"
	var buf bytes.Buffer
	buf.WriteByte(byte(data.Type))
	buf.WriteString(strconv.Itoa(int(data.Integer)))
	buf.Write(CRLF)
	b, err := w.writer.Write(buf.Bytes())
	return w.handleErrOnWrite(err, b, "interger", op)
}

func (w *Writer) simpleError(data *Value) error {
	op := "writer.simpleError"
	var buf bytes.Buffer
	buf.WriteByte(byte(data.Type))
	buf.Write(data.Bytes)
	buf.Write(CRLF)
	b, err := w.writer.Write(buf.Bytes())
	return w.handleErrOnWrite(err, b, "simpleError", op)
}

func (w *Writer) handleErrOnWrite(err error, b int, typ string, op string) error {
	if err != nil {

		w.logger.Error(op, slog.Any(ErrToWrite.Error(), err))
		return err
	}

	w.logger.Info(op, slog.Int(typ+" write size:", b))
	return nil
}
