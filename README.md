# golog

A simple logging library for Go, inspired by Log4j. It supports multiple log levels, console and file output, and log rotation.

## Installation

```bash
go get github.com/<your-username>/golog
```

## Usage

```go
package main

import (
	"github.com/<your-username>/golog"
)

func main() {
	logger, err := golog.NewLogger(golog.Config{
		Level:        golog.INFO,
		FilePath:     "app.log",
		LogToConsole: true,
	})
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Debug("This won't be logged")
	logger.Info("This is an info message")
	logger.Warn("This is a warning")
	logger.Error("This is an error")
	logger.Rotate() // Rotate log file
}
```

## Log Levels

- `DEBUG`: Detailed information for debugging.
- `INFO`: General information about application progress.
- `WARN`: Potentially harmful situations.
- `ERROR`: Errors that allow the application to continue.
- `FATAL`: Severe errors causing the application to exit.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to contribute.

## License

MIT License. See [LICENSE](LICENSE) for details.