// Package logging provides slog loggers with a request ID.
package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// New returns a logger with the request ID attached.
func New(handler slog.Handler, requestID string) *slog.Logger {
	logger := slog.New(handler)

	if requestID != "" {
		logger = logger.With(
			slog.String("request_id", requestID),
		)
	}

	return logger
}

// StdoutHandler returns a JSON handler for standard output.
func StdoutHandler(level slog.Level) slog.Handler {
	return slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)
}

// FileHandler returns a JSON handler for the file at path.
func FileHandler(
	path string,
	level slog.Level,
) (slog.Handler, io.Closer, error) {
	file, err := os.OpenFile(
		filepath.Clean(path),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return nil, nil, err
	}

	handler := slog.NewJSONHandler(
		file,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	return handler, file, nil
}
