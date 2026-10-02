package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func decodeRecords(t *testing.T, data []byte) []map[string]any {
	t.Helper()

	var records []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("unmarshal %q: %v", line, err)
		}
		records = append(records, record)
	}

	return records
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		requestID string
		want      any
		wantFound bool
	}{
		{name: "with request ID", requestID: "req-1", want: "req-1", wantFound: true},
		{name: "without request ID", requestID: "", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(slog.NewJSONHandler(&buf, nil), tt.requestID)

			logger.InfoContext(context.Background(), "hello")

			records := decodeRecords(t, buf.Bytes())
			if len(records) != 1 {
				t.Fatalf("got %d records, want 1", len(records))
			}
			if records[0]["msg"] != "hello" {
				t.Errorf("msg = %v, want hello", records[0]["msg"])
			}
			got, found := records[0]["request_id"]
			if found != tt.wantFound {
				t.Fatalf("request_id found = %v, want %v", found, tt.wantFound)
			}
			if found && got != tt.want {
				t.Errorf("request_id = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdoutHandler(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = stdout })

	logger := New(StdoutHandler(slog.LevelInfo), "req-1")
	logger.DebugContext(context.Background(), "debug")
	logger.InfoContext(context.Background(), "info")

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	records := decodeRecords(t, data)
	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}
	if records[0]["msg"] != "info" {
		t.Errorf("msg = %v, want info", records[0]["msg"])
	}
	if records[0]["request_id"] != "req-1" {
		t.Errorf("request_id = %v, want req-1", records[0]["request_id"])
	}
}

func TestFileHandler(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")

	for _, msg := range []string{"first", "second"} {
		handler, closer, err := FileHandler(path, slog.LevelInfo)
		if err != nil {
			t.Fatal(err)
		}
		logger := New(handler, "req-1")
		logger.DebugContext(context.Background(), "debug")
		logger.InfoContext(context.Background(), msg)
		if err := closer.Close(); err != nil {
			t.Fatal(err)
		}
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}

	records := decodeRecords(t, data)
	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}
	for i, want := range []string{"first", "second"} {
		if records[i]["msg"] != want {
			t.Errorf("records[%d].msg = %v, want %s", i, records[i]["msg"], want)
		}
		if records[i]["request_id"] != "req-1" {
			t.Errorf("records[%d].request_id = %v, want req-1", i, records[i]["request_id"])
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("perm = %o, want 600", perm)
	}
}

func TestFileHandlerError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "app.log")

	handler, closer, err := FileHandler(path, slog.LevelInfo)
	if err == nil {
		t.Fatal("err = nil, want error")
	}
	if handler != nil || closer != nil {
		t.Errorf("handler = %v, closer = %v, want nil", handler, closer)
	}
}
