package golog

import (
	"encoding/json"
	"fmt"
	"time"
)

// Formatter defines the interface for log formatting.
type Formatter interface {
	Format(level LogLevel, msg string, fields map[string]interface{}) string
}

// TextFormatter formats logs in plain text.
type TextFormatter struct{}

// Format implements text formatting.
func (f *TextFormatter) Format(level LogLevel, msg string, fields map[string]interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	base := fmt.Sprintf("[%s] %s %s", timestamp, level.String(), msg)
	if len(fields) == 0 {
		return base + "\n"
	}
	return fmt.Sprintf("%s %v\n", base, fields)
}

// JSONFormatter formats logs in JSON.
type JSONFormatter struct{}

// Format implements JSON formatting.
func (f *JSONFormatter) Format(level LogLevel, msg string, fields map[string]interface{}) string {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level.String(),
		"message":   msg,
	}
	for k, v := range fields {
		logEntry[k] = v
	}
	data, _ := json.Marshal(logEntry)
	return string(data) + "\n"
}
