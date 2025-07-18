package golog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerTextOutput(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger, err := NewLogger(Config{
		Level:        TRACE,
		FilePath:     logFile,
		LogToConsole: false,
		Format:       "text",
		MaxSizeMB:    1,
		MaxBackups:   3,
		Compress:     false,
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Test message", map[string]interface{}{"key": "value"})

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "INFO Test message map[key:value]") {
		t.Errorf("Log file does not contain expected message")
	}
}

func TestLoggerJSONOutput(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger, err := NewLogger(Config{
		Level:        TRACE,
		FilePath:     logFile,
		LogToConsole: false,
		Format:       "json",
		MaxSizeMB:    1,
		MaxBackups:   3,
		Compress:     false,
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Test JSON message", map[string]interface{}{"key": "value"})

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), `"message":"Test JSON message"`) || !strings.Contains(string(content), `"key":"value"`) {
		t.Errorf("Log file does not contain expected JSON message")
	}
}

func TestLogRotation(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger, err := NewLogger(Config{
		Level:        TRACE,
		FilePath:     logFile,
		LogToConsole: false,
		Format:       "text",
		MaxSizeMB:    1,
		MaxBackups:   2,
		Compress:     false,
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write enough data to trigger rotation
	for i := 0; i < 1000; i++ {
		logger.Info(strings.Repeat("x", 1024))
	}

	files, err := filepath.Glob(logFile + ".*")
	if err != nil {
		t.Fatalf("Failed to list log files: %v", err)
	}

	if len(files) < 1 {
		t.Errorf("Expected at least one rotated log file, got %d", len(files))
	}
}

func TestLogLevels(t *testing.T) {
	logger, err := NewLogger(Config{
		Level:        WARN,
		LogToConsole: false,
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Trace("This should not be logged")
	logger.Debug("This should not be logged")
	logger.Info("This should not be logged")
	logger.Warn("This should be logged")
	logger.Error("This should be logged")
	// Fatal is not tested as it exits the program
}
