package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewUsesPrettyLogsForDevelopmentAuto(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	logger, err := New(Options{
		Environment: "development",
		Format:      "auto",
		Output:      &output,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	logger.Info("hello", "answer", 42)

	line := strings.TrimSpace(output.String())
	if line == "" {
		t.Fatal("expected log output")
	}

	if json.Valid([]byte(line)) {
		t.Fatalf("expected pretty output, got JSON: %s", line)
	}
}

func TestNewUsesJSONLogsForProductionAuto(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	logger, err := New(Options{
		Environment: "production",
		Format:      "auto",
		Output:      &output,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	logger.Info("hello", "answer", 42)

	record := decodeLogLine(t, output.String())
	if record["msg"] != "hello" {
		t.Fatalf("expected message hello, got %#v", record["msg"])
	}

	if record["service"] != "backend" {
		t.Fatalf("expected service backend, got %#v", record["service"])
	}

	if record["env"] != "production" {
		t.Fatalf("expected env production, got %#v", record["env"])
	}
}

func TestNewAllowsExplicitFormatOverride(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	logger, err := New(Options{
		Environment: "development",
		Format:      "json",
		Output:      &output,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	logger.Info("hello")
	if !json.Valid(bytes.TrimSpace(output.Bytes())) {
		t.Fatalf("expected JSON output, got %q", output.String())
	}
}

func TestFromContextReturnsStoredLogger(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	ctx := WithContext(context.Background(), logger)

	if got := FromContext(ctx); got != logger {
		t.Fatal("expected logger from context")
	}
}

func TestFromContextFallsBackToDefaultLogger(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	previous := slog.Default()
	slog.SetDefault(logger)
	t.Cleanup(func() {
		slog.SetDefault(previous)
	})

	got := FromContext(context.Background())
	if got != logger {
		t.Fatal("expected default logger fallback")
	}
}

func decodeLogLine(t *testing.T, raw string) map[string]any {
	t.Helper()

	line := strings.TrimSpace(raw)
	if line == "" {
		t.Fatal("expected log output")
	}

	var record map[string]any
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		t.Fatalf("unmarshal log line: %v", err)
	}

	return record
}
