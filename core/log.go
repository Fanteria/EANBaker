package core

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type MultiLogger struct {
	*slog.Logger
	buffer *syncBuffer
	multi  io.Writer
}

// thread-safe buffer
type syncBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (b *syncBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Bytes()
}

func (b *syncBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.Reset()
}

func MultiLoggerFromEnv() (*MultiLogger, error) {
	env_log := os.Getenv("LOG")
	switch strings.ToLower(env_log) {
	case "error":
		return NewMultiLogger(slog.LevelError), nil
	case "info", "", "default":
		return NewMultiLogger(slog.LevelInfo), nil
	case "debug":
		return NewMultiLogger(slog.LevelDebug), nil
	default:
		return nil, fmt.Errorf("Invalid logging level %s", env_log)
	}
}

// New creates a new logger that writes to terminal and keeps logs in memory
func NewMultiLogger(level slog.Level) *MultiLogger {
	buf := &syncBuffer{}

	// Write to both terminal and buffer
	multi := io.MultiWriter(os.Stderr, buf)

	handler := slog.NewJSONHandler(multi, &slog.HandlerOptions{
		Level: level,
	})

	return &MultiLogger{
		Logger: slog.New(handler),
		buffer: buf,
		multi:  multi,
	}
}

// SaveToFile writes all buffered logs to a file
func (l *MultiLogger) SaveToFile(path string) error {
	return os.WriteFile(path, l.buffer.Bytes(), 0644)
}
