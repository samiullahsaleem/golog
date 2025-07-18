package golog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger, err := NewLogger(Config{
		Level:        DEBUG,
		FilePath:     logFile,
		LogToConsole: false,
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Test message: %s", "hello")
	logger.Rotate()

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "INFO Test message: hello") {
		t.Errorf("Log file does not contain expected message")
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

	logger.Debug("This should not be logged")
	logger.Info("This should not be logged")
	logger.Warn("This should be logged")
	logger.Error("This should be logged")
	logger.Fatal("This should not run") // Fatal exits, so avoid in test
}
