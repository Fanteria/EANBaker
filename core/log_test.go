package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"log/slog"
)

func TestSyncBuffer_Write(t *testing.T) {
	buf := &syncBuffer{}

	data := []byte("test data")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() wrote %d bytes, want %d", n, len(data))
	}
}

func TestSyncBuffer_Bytes(t *testing.T) {
	buf := &syncBuffer{}

	_, err := buf.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	bytes := buf.Bytes()
	if string(bytes) != "test" {
		t.Errorf("Bytes() = %v, want %v", string(bytes), "test")
	}
}

func TestSyncBuffer_Reset(t *testing.T) {
	buf := &syncBuffer{}

	_, _ = buf.Write([]byte("test"))
	buf.Reset()

	bytes := buf.Bytes()
	if len(bytes) != 0 {
		t.Errorf("Reset() didn't clear buffer, got %d bytes", len(bytes))
	}
}

func TestSyncBuffer_MultipleWrites(t *testing.T) {
	buf := &syncBuffer{}

	_, _ = buf.Write([]byte("hello"))
	_, _ = buf.Write([]byte(" "))
	_, _ = buf.Write([]byte("world"))

	result := string(buf.Bytes())
	if result != "hello world" {
		t.Errorf("Multiple writes = %v, want %v", result, "hello world")
	}
}

func TestNewMultiLogger(t *testing.T) {
	tests := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}

	for _, level := range tests {
		t.Run(level.String(), func(t *testing.T) {
			logger := NewMultiLogger(level)
			if logger == nil {
				t.Fatal("NewMultiLogger() returned nil")
			}
			if logger.Logger == nil {
				t.Error("NewMultiLogger().Logger is nil")
			}
			if logger.buffer == nil {
				t.Error("NewMultiLogger().buffer is nil")
			}
			if logger.multi == nil {
				t.Error("NewMultiLogger().multi is nil")
			}
		})
	}
}

func TestMultiLogger_Logging(t *testing.T) {
	logger := NewMultiLogger(slog.LevelDebug)

	// Log some messages
	logger.Info("info message", "key", "value")
	logger.Debug("debug message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Check buffer contains logs
	logs := string(logger.buffer.Bytes())
	if !strings.Contains(logs, "info message") {
		t.Error("Buffer should contain 'info message'")
	}
	if !strings.Contains(logs, "debug message") {
		t.Error("Buffer should contain 'debug message'")
	}
	if !strings.Contains(logs, "warn message") {
		t.Error("Buffer should contain 'warn message'")
	}
	if !strings.Contains(logs, "error message") {
		t.Error("Buffer should contain 'error message'")
	}
}

func TestMultiLogger_LevelFiltering(t *testing.T) {
	// Create logger with INFO level - should filter out DEBUG
	logger := NewMultiLogger(slog.LevelInfo)

	logger.Debug("debug message")
	logger.Info("info message")

	logs := string(logger.buffer.Bytes())
	if strings.Contains(logs, "debug message") {
		t.Error("DEBUG messages should be filtered at INFO level")
	}
	if !strings.Contains(logs, "info message") {
		t.Error("INFO messages should be logged at INFO level")
	}
}

func TestMultiLogger_ErrorLevel(t *testing.T) {
	// Create logger with ERROR level
	logger := NewMultiLogger(slog.LevelError)

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	logs := string(logger.buffer.Bytes())
	if strings.Contains(logs, "debug") {
		t.Error("DEBUG should be filtered")
	}
	if strings.Contains(logs, "info") {
		t.Error("INFO should be filtered")
	}
	if strings.Contains(logs, "warn") {
		t.Error("WARN should be filtered")
	}
	if !strings.Contains(logs, "error") {
		t.Error("ERROR should be logged")
	}
}

func TestMultiLogger_SaveToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "log-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logger := NewMultiLogger(slog.LevelInfo)
	logger.Info("test message", "key", "value")

	logPath := filepath.Join(tmpDir, "test.log")
	err = logger.SaveToFile(logPath)
	if err != nil {
		t.Errorf("SaveToFile() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("Log file was not created")
	}

	// Verify content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "test message") {
		t.Error("Log file should contain 'test message'")
	}
}

func TestMultiLogger_SaveToFile_InvalidPath(t *testing.T) {
	logger := NewMultiLogger(slog.LevelInfo)
	logger.Info("test")

	err := logger.SaveToFile("/nonexistent/directory/test.log")
	if err == nil {
		t.Error("SaveToFile() should fail for invalid path")
	}
}

func TestMultiLogger_SaveToFile_EmptyBuffer(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "log-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logger := NewMultiLogger(slog.LevelInfo)
	// No logging done

	logPath := filepath.Join(tmpDir, "empty.log")
	err = logger.SaveToFile(logPath)
	if err != nil {
		t.Errorf("SaveToFile() failed: %v", err)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if len(content) != 0 {
		t.Errorf("Empty log file should be empty, got %d bytes", len(content))
	}
}

func TestMultiLoggerFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envVal  string
		wantErr bool
	}{
		{name: "Empty env", envVal: "", wantErr: false},
		{name: "Default", envVal: "default", wantErr: false},
		{name: "Info lowercase", envVal: "info", wantErr: false},
		{name: "Info uppercase", envVal: "INFO", wantErr: false},
		{name: "Info mixed case", envVal: "Info", wantErr: false},
		{name: "Debug lowercase", envVal: "debug", wantErr: false},
		{name: "Debug uppercase", envVal: "DEBUG", wantErr: false},
		{name: "Error lowercase", envVal: "error", wantErr: false},
		{name: "Error uppercase", envVal: "ERROR", wantErr: false},
		{name: "Invalid value", envVal: "invalid", wantErr: true},
		{name: "Numeric value", envVal: "123", wantErr: true},
		{name: "Warn (not supported)", envVal: "warn", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env variable
			originalVal := os.Getenv("LOG")
			os.Setenv("LOG", tt.envVal)
			defer os.Setenv("LOG", originalVal)

			logger, err := MultiLoggerFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiLoggerFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("MultiLoggerFromEnv() returned nil logger")
			}
		})
	}
}

func TestMultiLogger_JSONFormat(t *testing.T) {
	logger := NewMultiLogger(slog.LevelInfo)
	logger.Info("test message", "key", "value")

	logs := string(logger.buffer.Bytes())

	// JSON format should contain these elements
	if !strings.Contains(logs, `"msg"`) {
		t.Error("JSON output should contain 'msg' field")
	}
	if !strings.Contains(logs, `"level"`) {
		t.Error("JSON output should contain 'level' field")
	}
	if !strings.Contains(logs, `"time"`) {
		t.Error("JSON output should contain 'time' field")
	}
}

func TestMultiLogger_JSONAttributes(t *testing.T) {
	logger := NewMultiLogger(slog.LevelInfo)
	logger.Info("message", "string_key", "string_value", "int_key", 42, "bool_key", true)

	logs := string(logger.buffer.Bytes())

	if !strings.Contains(logs, `"string_key"`) {
		t.Error("JSON should contain string_key")
	}
	if !strings.Contains(logs, `"int_key"`) {
		t.Error("JSON should contain int_key")
	}
	if !strings.Contains(logs, `"bool_key"`) {
		t.Error("JSON should contain bool_key")
	}
}

func TestSyncBuffer_ConcurrentWrite(t *testing.T) {
	buf := &syncBuffer{}

	// Write from multiple goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				_, _ = buf.Write([]byte("x"))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 1000 'x' characters
	if len(buf.Bytes()) != 1000 {
		t.Errorf("Expected 1000 bytes, got %d", len(buf.Bytes()))
	}
}

func TestMultiLogger_ConcurrentLogging(t *testing.T) {
	logger := NewMultiLogger(slog.LevelInfo)

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 50; j++ {
				logger.Info("message", "goroutine", id, "iteration", j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have logged 500 messages
	logs := string(logger.buffer.Bytes())
	count := strings.Count(logs, `"msg"`)
	if count != 500 {
		t.Errorf("Expected 500 log entries, got %d", count)
	}
}
