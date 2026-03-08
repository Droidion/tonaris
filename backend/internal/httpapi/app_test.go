package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if body := recorder.Body.String(); body != "Hello World" {
		t.Fatalf("expected body %q, got %q", "Hello World", body)
	}

	if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
		t.Fatalf("expected text/plain content type, got %q", contentType)
	}
}

func TestOpenAPIEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	request := httptest.NewRequest(http.MethodGet, "/api-doc/openapi.json", nil)
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"/hello"`) {
		t.Fatalf("expected openapi document to include /hello path, got %q", recorder.Body.String())
	}
}

func TestDocsEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	request := httptest.NewRequest(http.MethodGet, "/scalar", nil)
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code < http.StatusOK || recorder.Code >= http.StatusBadRequest {
		t.Fatalf("expected docs status below 400, got %d", recorder.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	request := httptest.NewRequest(http.MethodOptions, "/hello", nil)
	request.Header.Set("Origin", "https://example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected wildcard allow origin, got %q", recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestOpenAPIDocumentMatchesExpectedShape(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	document := app.API.OpenAPI()

	if document.Info == nil {
		t.Fatal("expected document info")
	}

	if document.Info.Title != apiTitle {
		t.Fatalf("expected title %q, got %q", apiTitle, document.Info.Title)
	}

	if document.Info.Version != apiVersion {
		t.Fatalf("expected version %q, got %q", apiVersion, document.Info.Version)
	}

	operation := document.Paths["/hello"].Get
	if operation == nil {
		t.Fatal("expected /hello GET operation")
	}

	if operation.OperationID != "tonaris.hello" {
		t.Fatalf("expected operationId tonaris.hello, got %q", operation.OperationID)
	}

	response := operation.Responses["200"]
	if response == nil {
		t.Fatal("expected 200 response")
	}

	content := response.Content["text/plain"]
	if content == nil {
		t.Fatal("expected text/plain response content")
	}

	if content.Schema == nil || content.Schema.Type != "string" {
		t.Fatalf("expected string schema, got %#v", content.Schema)
	}
}

func newTestApp(t *testing.T) *App {
	t.Helper()

	app, err := New()
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	return app
}
