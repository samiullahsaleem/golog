package golog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// Logger represents a logging instance.
type Logger struct {
	level        LogLevel
	file         *os.File
	filePath     string
	mutex        sync.Mutex
	logToFile    bool
	logToConsole bool
}

// Config holds logger configuration options.
type Config struct {
	Level        LogLevel
	FilePath     string
	LogToConsole bool
}

// NewLogger creates a new logger with the given configuration.
func NewLogger(config Config) (*Logger, error) {
	logger := &Logger{
		level:        config.Level,
		logToFile:    config.FilePath != "",
		logToConsole: config.LogToConsole,
	}

	if logger.logToFile {
		err := logger.openLogFile(config.FilePath)
		if err != nil {
			return nil, err
		}
	}

	return logger, nil
}

// openLogFile opens or creates the log file.
func (l *Logger) openLogFile(filePath string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	l.filePath = filePath
	l.file = file
	return nil
}

// log writes a log message if the level is sufficient.
func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := fmt.Sprintf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), level.String(), fmt.Sprintf(msg, args...))

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.logToConsole {
		fmt.Print(message)
	}

	if l.logToFile && l.file != nil {
		l.file.WriteString(message)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

// Info logs an info message.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

// Fatal logs a fatal message and exits the program.
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

// Close closes the log file.
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Rotate rotates the log file, renaming the current file with a timestamp.
func (l *Logger) Rotate() error {
	if !l.logToFile || l.file == nil {
		return nil
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if err := l.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %v", err)
	}

	newPath := fmt.Sprintf("%s.%s", l.filePath, time.Now().Format("20060102_150405"))
	if err := os.Rename(l.filePath, newPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %v", err)
	}

	return l.openLogFile(l.filePath)
}
