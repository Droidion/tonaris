package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func TestMiddlewareLogsSuccessfulRequest(t *testing.T) {
	t.Parallel()

	recorder, entries := serveRequest(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := FromContext(r.Context())
		logger.InfoContext(r.Context(), "inside handler", "extra", "value")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	}))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}

	completion := findEntry(t, entries, "request completed")
	if completion["level"] != "INFO" {
		t.Fatalf("expected INFO level, got %#v", completion["level"])
	}

	if completion["status"] != float64(http.StatusCreated) {
		t.Fatalf("expected status 201, got %#v", completion["status"])
	}

	if completion["method"] != http.MethodGet {
		t.Fatalf("expected GET method, got %#v", completion["method"])
	}

	if completion["path"] != "/hello" {
		t.Fatalf("expected path /hello, got %#v", completion["path"])
	}

	if completion["remote_ip"] != "203.0.113.10" {
		t.Fatalf("expected remote ip 203.0.113.10, got %#v", completion["remote_ip"])
	}

	if completion["user_agent"] != "tonaris-test" {
		t.Fatalf("expected user agent tonaris-test, got %#v", completion["user_agent"])
	}

	if _, ok := completion["request_id"]; !ok {
		t.Fatal("expected request_id field")
	}

	if _, ok := completion["duration_ms"]; !ok {
		t.Fatal("expected duration_ms field")
	}

	if completion["bytes"] != float64(2) {
		t.Fatalf("expected bytes 2, got %#v", completion["bytes"])
	}
}

func TestMiddlewareLogsWarnForClientError(t *testing.T) {
	t.Parallel()

	_, entries := serveRequest(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))

	completion := findEntry(t, entries, "request completed")
	if completion["level"] != "WARN" {
		t.Fatalf("expected WARN level, got %#v", completion["level"])
	}
}

func TestMiddlewareLogsErrorForServerError(t *testing.T) {
	t.Parallel()

	_, entries := serveRequest(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))

	completion := findEntry(t, entries, "request completed")
	if completion["level"] != "ERROR" {
		t.Fatalf("expected ERROR level, got %#v", completion["level"])
	}
}

func TestMiddlewareRecoversPanics(t *testing.T) {
	t.Parallel()

	recorder, entries := serveRequest(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("kaboom")
	}))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}

	panicEntry := findEntry(t, entries, "panic recovered")
	if panicEntry["level"] != "ERROR" {
		t.Fatalf("expected ERROR level, got %#v", panicEntry["level"])
	}

	if panicEntry["panic"] != "kaboom" {
		t.Fatalf("expected panic payload kaboom, got %#v", panicEntry["panic"])
	}

	stackTrace, ok := panicEntry["stack_trace"].(string)
	if !ok || stackTrace == "" {
		t.Fatal("expected stack_trace field")
	}

	completion := findEntry(t, entries, "request completed")
	if completion["status"] != float64(http.StatusInternalServerError) {
		t.Fatalf("expected completion status 500, got %#v", completion["status"])
	}
}

func TestMiddlewareMakesLoggerAvailableInContext(t *testing.T) {
	t.Parallel()

	errLoggerMissing := errors.New("logger missing from context")
	var sawRequestLogger bool

	_, entries := serveRequest(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := FromContext(r.Context())
		if logger == slog.Default() {
			panic(errLoggerMissing)
		}

		sawRequestLogger = true
		logger.InfoContext(r.Context(), "context logger works")
		w.WriteHeader(http.StatusNoContent)
	}))

	if !sawRequestLogger {
		t.Fatal("expected handler to observe request logger")
	}

	findEntry(t, entries, "context logger works")
}

func serveRequest(t *testing.T, next http.Handler) (*httptest.ResponseRecorder, []map[string]any) {
	t.Helper()

	var output bytes.Buffer
	logger, err := New(Options{
		Environment: "production",
		Format:      "json",
		Output:      &output,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	handler := chimiddleware.RequestID(chimiddleware.RealIP(Middleware(logger)(next)))
	request := httptest.NewRequest(http.MethodGet, "/hello?secret=1", nil)
	request = request.WithContext(context.Background())
	request.RemoteAddr = "198.51.100.7:1234"
	request.Header.Set("X-Forwarded-For", "203.0.113.10")
	request.Header.Set("User-Agent", "tonaris-test")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder, decodeLogEntries(t, output.String())
}

func decodeLogEntries(t *testing.T, raw string) []map[string]any {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(raw), "\n")
	entries := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("unmarshal log entry: %v", err)
		}

		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		t.Fatal("expected log entries")
	}

	return entries
}

func findEntry(t *testing.T, entries []map[string]any, message string) map[string]any {
	t.Helper()

	for _, entry := range entries {
		if entry["msg"] == message {
			return entry
		}
	}

	t.Fatalf("expected log entry %q", message)
	return nil
}
