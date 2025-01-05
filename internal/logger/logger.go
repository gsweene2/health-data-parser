// Package logger provides a configured zap logger for consistent logging across the application.
package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// New returns a configured SugaredLogger instance. The logger is configured in
// production mode. If initialization fails, the program exits with status code 1.
func New() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	return logger.Sugar()
}
