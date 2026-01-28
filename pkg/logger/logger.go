package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	once   sync.Once
	logger *slog.Logger
)

func InitLogger(service string, logLevel string, logFile string) *slog.Logger {
	once.Do(func() {
		level := getLoggerLevel(logLevel)
		handler := createHandler(level, service, logFile)

		logger = slog.New(handler)
	})
	return logger
}

func getLoggerLevel(logLevel string) slog.Level {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		logLevel = envLogLevel
	}
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type CustomHandler struct {
	Level   slog.Level
	Service string
	Writer  io.Writer
}

func (c *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= c.Level
}

func (c *CustomHandler) Handle(_ context.Context, record slog.Record) error {
	timestamp := record.Time.UTC().Format(time.RFC3339Nano)
	level := strings.ToUpper(record.Level.String())

	message := record.Message

	var attrs []string
	record.Attrs(func(a slog.Attr) bool {
		if a.Key == "!BADKEY" {
			a.Key = ""
		}
		attrs = append(attrs, a.Key)

		attrs = append(attrs, a.Value.String())
		return true
	})
	var extra string
	if len(attrs) > 0 {
		extra = strings.Join(attrs, " ")
	}

	_, err := fmt.Fprintf(c.Writer, "%s [%s] [%s] %s %s\n", timestamp, level, c.Service, message, extra)
	return err
}

func (c *CustomHandler) WithAttrs([]slog.Attr) slog.Handler {
	return c
}

func (c *CustomHandler) WithGroup(string) slog.Handler {
	return c
}

func createHandler(level slog.Level, service string, logFile string) slog.Handler {
	if logFile == "" {
		logFile = os.Getenv("LOG_FILE")
	}
	dir := filepath.Dir(logFile)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic("could not open log folder: " + err.Error())
	}
	file, errF := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if errF != nil {
		panic("could not open log file: " + errF.Error())
	}

	writer := io.MultiWriter(os.Stdout, file)

	return &CustomHandler{
		Level:   level,
		Service: service,
		Writer:  writer,
	}
}

func GetLogger() *slog.Logger {
	if logger == nil {
		panic("Logger is not initialized. Call InitLogger first.")
	}
	return logger
}

func LoggerNotInitialized(logger *slog.Logger) *slog.Logger {
	op := "logger.LoggerInitialized"
	if logger == nil {
		logger = slog.New(
			slog.NewTextHandler(io.Discard, nil),
		)
		log.Print(op + " no logger provided")
	}
	return logger
}
