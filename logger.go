package golog

import (
	"fmt"
	"os"
	"sync"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	return [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// Logger represents a logging instance.
type Logger struct {
	level        LogLevel
	formatter    Formatter
	file         *os.File
	filePath     string
	mutex        sync.Mutex
	logToFile    bool
	logToConsole bool
	rotator      *Rotator
}

// Config holds logger configuration options.
type Config struct {
	Level        LogLevel
	FilePath     string
	LogToConsole bool
	Format       string // "text" or "json"
	MaxSizeMB    int    // Max file size in MB before rotation
	MaxBackups   int    // Max number of backup files
	Compress     bool   // Compress rotated files
}

// NewLogger creates a new logger with the given configuration.
func NewLogger(config Config) (*Logger, error) {
	logger := &Logger{
		level:        config.Level,
		logToFile:    config.FilePath != "",
		logToConsole: config.LogToConsole,
	}

	if config.Format == "json" {
		logger.formatter = &JSONFormatter{}
	} else {
		logger.formatter = &TextFormatter{}
	}

	if logger.logToFile {
		var err error
		logger.file, err = os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		logger.filePath = config.FilePath
		logger.rotator = NewRotator(config.FilePath, config.MaxSizeMB, config.MaxBackups, config.Compress)
	}

	return logger, nil
}

// log writes a log message if the level is sufficient.
func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	message := l.formatter.Format(level, msg, fields)

	if l.logToConsole {
		fmt.Print(message)
	}

	if l.logToFile && l.file != nil {
		if l.rotator != nil {
			if err := l.rotator.RotateIfNeeded(l.file); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to rotate log: %v\n", err)
			}
		}
		l.file.WriteString(message)
	}
}

// Trace logs a trace message.
func (l *Logger) Trace(msg string, fields ...map[string]interface{}) {
	l.log(TRACE, msg, mergeFields(fields))
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(DEBUG, msg, mergeFields(fields))
}

// Info logs an info message.
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	l.log(INFO, msg, mergeFields(fields))
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(WARN, msg, mergeFields(fields))
}

// Error logs an error message.
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	l.log(ERROR, msg, mergeFields(fields))
}

// Fatal logs a fatal message and exits the program.
func (l *Logger) Fatal(msg string, fields ...map[string]interface{}) {
	l.log(FATAL, msg, mergeFields(fields))
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

// mergeFields combines multiple field maps into one.
func mergeFields(fields []map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, f := range fields {
		for k, v := range f {
			result[k] = v
		}
	}
	return result
}
