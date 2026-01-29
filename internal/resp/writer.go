package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"

	l "bedis/pkg/logger"
)

var (
	ErrToWrite  = errors.New("unable to write in writer")
	ErrNilValue = errors.New("value is undefined")
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

func (w *Writer) WriteValue(value *Value) error {
	op := "writer.WriterValue"
	if value == nil {
		w.logger.Error(op, slog.Any("ERROR:", ErrNilValue))
		return ErrNilValue
	}
	switch value.Type {
	case SimpleString:
		return w.writeSimpleString(value)
	default:
		return ErrUnknownType
	}
}

// func (w *Writer) writeSimpleError(data *Value) {
// 	op := "writer.writeSimpleError"
// }

func (w *Writer) writeSimpleString(data *Value) error {
	op := "writer.SimpleString"
	var buf bytes.Buffer
	buf.WriteByte(byte(data.Type))
	buf.Write(data.Bytes)
	buf.Write([]byte{'\r', '\n'})
	b, err := w.writer.Write(buf.Bytes())
	if err != nil {
		w.logger.Error(op, slog.Any(ErrToWrite.Error(), err))
		return err
	}

	w.logger.Info(op, slog.Int("SimpleString write size:", b))
	return nil
}
