# golog

`golog` is a lightweight, thread-safe logging library for Go, inspired by Log4j. It provides flexible logging with multiple log levels, structured logging, customizable output formats (text or JSON), and robust log rotation with optional compression. Designed for simplicity and performance, `golog` is ideal for both development and production environments.

## Features

- **Log Levels**: TRACE, DEBUG, INFO, WARN, ERROR, FATAL
- **Output Formats**: Plain text or JSON for structured logging
- **Structured Logging**: Attach key-value pairs to logs for better context
- **Log Rotation**: Size-based rotation with configurable maximum file size and backup count
- **Compression**: Optional gzip compression for rotated log files
- **Thread-Safe**: Safe for concurrent use in multi-goroutine applications
- **Configurable Outputs**: Log to console, file, or both
- **Go 1.24.5 Compatible**: Tested and optimized for the latest Go version

## Installation

To install `golog`, use the following command:

```bash
go get github.com/samiullahsaleem/golog@v1.0.0
```

Ensure you are using Go 1.24.5 or later. Verify your Go version with:

```bash
go version
```

## Usage

Below is an example of how to use `golog` in your Go application:

```go
package main

import (
	"github.com/samiullahsaleem/golog"
)

func main() {
	// Configure the logger
	logger, err := golog.NewLogger(golog.Config{
		Level:        golog.INFO,      // Log INFO and above
		FilePath:     "app.log",       // Output to app.log
		LogToConsole: true,            // Also output to console
		Format:       "json",          // Use JSON format
		MaxSizeMB:    10,              // Rotate after 10MB
		MaxBackups:   5,               // Keep up to 5 backups
		Compress:     true,            // Compress rotated files
	})
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// Log messages with different levels and structured data
	logger.Trace("Trace message", map[string]interface{}{"user": "alice"})
	logger.Debug("Debug message", map[string]interface{}{"action": "start"})
	logger.Info("Application started", map[string]interface{}{"version": "1.0.0"})
	logger.Warn("Low memory", map[string]interface{}{"memory_mb": 100})
	logger.Error("Connection failed", map[string]interface{}{"error": "timeout"})

	// Simulate large log output to trigger rotation
	for i := 0; i < 1000; i++ {
		logger.Info("Large log data", map[string]interface{}{"index": i})
	}
}
```

### Example Output

- **Console (JSON format)**:
  ```json
  {"timestamp":"2025-07-18T21:48:00Z","level":"INFO","message":"Application started","version":"1.0.0"}
  {"timestamp":"2025-07-18T21:48:00Z","level":"WARN","message":"Low memory","memory_mb":100}
  {"timestamp":"2025-07-18T21:48:00Z","level":"ERROR","message":"Connection failed","error":"timeout"}
  ...
  ```

- **File (`app.log`)**: Same content as console, with rotated files (e.g., `app.log.20250718_214800.gz`) created when the file exceeds 10MB.

## Log Levels

`golog` supports the following log levels:

- `TRACE`: Fine-grained diagnostic information for detailed debugging.
- `DEBUG`: Information useful during development and debugging.
- `INFO`: General information about application progress.
- `WARN`: Indications of potential issues that don’t halt execution.
- `ERROR`: Errors that allow the application to continue running.
- `FATAL`: Critical errors that cause the application to exit.

Set the desired level in the `Config.Level` field to filter logs. For example, setting `Level: golog.WARN` will only log WARN, ERROR, and FATAL messages.

## Configuration Options

The `golog.Config` struct allows you to customize the logger:

- `Level`: Minimum log level to record (e.g., `golog.INFO`).
- `FilePath`: Path to the log file (e.g., `"app.log"`). Set to empty string to disable file output.
- `LogToConsole`: Enable/disable console output (`true`/`false`).
- `Format`: Output format (`"text"` for plain text, `"json"` for structured JSON).
- `MaxSizeMB`: Maximum log file size in megabytes before rotation.
- `MaxBackups`: Maximum number of rotated log files to keep.
- `Compress`: Enable gzip compression for rotated log files.

## Log Rotation

`golog` automatically rotates log files when they exceed `MaxSizeMB`. Rotated files are named with a timestamp (e.g., `app.log.20250718_214800`). If `Compress` is `true`, rotated files are compressed with gzip (e.g., `app.log.20250718_214800.gz`). The `MaxBackups` setting limits the number of retained backups, deleting the oldest files when the limit is exceeded.

## Structured Logging

Attach key-value pairs to logs for additional context:

```go
logger.Info("User logged in", map[string]interface{}{"user_id": 123, "ip": "192.168.1.1"})
```

In JSON format, this produces:
```json
{"timestamp":"2025-07-18T21:48:00Z","level":"INFO","message":"User logged in","user_id":123,"ip":"192.168.1.1"}
```

In text format, it looks like:
```
[2025-07-18 21:48:00] INFO User logged in map[user_id:123 ip:192.168.1.1]
```

## Testing Locally

To test `golog` locally:

1. Clone the repository:
   ```bash
   git clone https://github.com/samiullahsaleem/golog.git
   cd golog
   ```

2. Run the unit tests:
   ```bash
   go test ./... -v
   ```

3. Create a sample program in a separate directory (e.g., `golog-example`):
   ```bash
   mkdir golog-example
   cd golog-example
   go mod init golog-example
   ```

4. Add `golog` as a dependency in `golog-example/go.mod`:
   ```go
   module golog-example

   go 1.24.5

   require github.com/samiullahsaleem/golog v1.0.0

   replace github.com/samiullahsaleem/golog => ../golog
   ```

5. Write a test program (see the Usage section) and run it:
   ```bash
   go run main.go
   ```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute, including code style and pull request processes.

## License

`golog` is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Contact

For issues, feature requests, or questions, please open an issue on the [GitHub repository](https://github.com/samiullahsaleem/golog).