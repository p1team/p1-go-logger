package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func New(handler slog.Handler, requestID string) *slog.Logger {
	logger := slog.New(handler)

	if requestID != "" {
		logger = logger.With(
			slog.String("request_id", requestID),
		)
	}

	return logger
}

func StdoutHandler(level slog.Level) slog.Handler {
	return slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)
}

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
