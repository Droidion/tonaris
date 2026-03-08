package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/apperr"

	"github.com/danielgtaylor/huma/v2"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
	recorder := performRequest(app, http.MethodGet, "/healthz")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var body healthStatus
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if body.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Status)
	}

	if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("expected application/json content type, got %q", contentType)
	}
}

func TestHelloEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
	recorder := performRequest(app, http.MethodGet, "/hello")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if body := recorder.Body.String(); body != "hello world" {
		t.Fatalf("expected body %q, got %q", "hello world", body)
	}

	if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
		t.Fatalf("expected text/plain content type, got %q", contentType)
	}
}

func TestOpenAPIEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
	recorder := performRequest(app, http.MethodGet, "/api-doc/openapi.json")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"/healthz"`) {
		t.Fatalf("expected openapi document to include /healthz path, got %q", recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), `"/hello"`) {
		t.Fatalf("expected openapi document to include /hello path, got %q", recorder.Body.String())
	}
}

func TestDocsEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
	recorder := performRequest(app, http.MethodGet, "/scalar")

	if recorder.Code < http.StatusOK || recorder.Code >= http.StatusBadRequest {
		t.Fatalf("expected docs status below 400, got %d", recorder.Code)
	}
}

func TestCORSPreflightUsesConfiguredOrigins(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{AllowedOrigins: []string{"https://example.com"}})
	request := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
	request.Header.Set("Origin", "https://example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Fatalf("expected configured allow origin, got %q", recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestAppUsesDefaultOriginsWhenDependenciesAreEmpty(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
	request := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("expected localhost dev origin, got %q", recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestOpenAPIDocumentMatchesExpectedShape(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
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

	operation := document.Paths["/healthz"].Get
	if operation == nil {
		t.Fatal("expected /healthz GET operation")
	}

	if operation.OperationID != "tonaris.system.health" {
		t.Fatalf("expected operationId tonaris.system.health, got %q", operation.OperationID)
	}

	helloOperation := document.Paths["/hello"].Get
	if helloOperation == nil {
		t.Fatal("expected /hello GET operation")
	}

	if helloOperation.OperationID != "tonaris.system.hello" {
		t.Fatalf("expected operationId tonaris.system.hello, got %q", helloOperation.OperationID)
	}
}

func TestAppErrorIsMappedToProblemResponse(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{}, func(api huma.API) {
		huma.Register(api, huma.Operation{
			OperationID: "test.not-found",
			Method:      http.MethodGet,
			Path:        "/test/not-found",
		}, func(context.Context, *struct{}) (*struct{}, error) {
			return nil, apperr.New(apperr.NotFound, "users.not_found", "user not found")
		})
	})

	recorder := performRequest(app, http.MethodGet, "/test/not-found")

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/problem+json") {
		t.Fatalf("expected problem json content type, got %q", contentType)
	}

	problem := decodeProblem(t, recorder)
	if problem.Code != "users.not_found" {
		t.Fatalf("expected users.not_found code, got %q", problem.Code)
	}

	if problem.Detail != "user not found" {
		t.Fatalf("expected safe detail, got %q", problem.Detail)
	}
}

func TestGenericErrorIsMappedToSafeInternalProblem(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{}, func(api huma.API) {
		huma.Register(api, huma.Operation{
			OperationID: "test.internal-error",
			Method:      http.MethodGet,
			Path:        "/test/internal-error",
		}, func(context.Context, *struct{}) (*struct{}, error) {
			return nil, errors.New("database connection failed")
		})
	})

	recorder := performRequest(app, http.MethodGet, "/test/internal-error")

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}

	problem := decodeProblem(t, recorder)
	if problem.Code != "internal.error" {
		t.Fatalf("expected internal.error code, got %q", problem.Code)
	}

	if problem.Detail != "unexpected error occurred" {
		t.Fatalf("expected safe internal detail, got %q", problem.Detail)
	}

	if len(problem.Errors) != 0 {
		t.Fatalf("expected no internal error details to be exposed, got %#v", problem.Errors)
	}
}

func TestValidationErrorKeepsStructuredErrorDetails(t *testing.T) {
	t.Parallel()

	type input struct {
		ID int `path:"id"`
	}

	app := newTestApp(t, Dependencies{}, func(api huma.API) {
		huma.Register(api, huma.Operation{
			OperationID: "test.validation",
			Method:      http.MethodGet,
			Path:        "/test/items/{id}",
		}, func(context.Context, *input) (*struct{}, error) {
			return &struct{}{}, nil
		})
	})

	recorder := performRequest(app, http.MethodGet, "/test/items/not-an-int")

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", recorder.Code)
	}

	problem := decodeProblem(t, recorder)
	if problem.Code != "http.unprocessable_entity" {
		t.Fatalf("expected http.unprocessable_entity code, got %q", problem.Code)
	}

	if problem.Detail != "validation failed" {
		t.Fatalf("expected validation failed detail, got %q", problem.Detail)
	}

	if len(problem.Errors) == 0 {
		t.Fatal("expected validation errors to be included")
	}
}

func TestOpenAPIDocumentUsesProblemSchema(t *testing.T) {
	t.Parallel()

	app := newTestApp(t, Dependencies{})
	document := app.API.OpenAPI()

	problemSchema := document.Components.Schemas.Map()["Problem"]
	if problemSchema == nil {
		t.Fatal("expected Problem schema in components")
	}

	for _, field := range []string{"code", "errors"} {
		if _, ok := problemSchema.Properties[field]; !ok {
			t.Fatalf("expected Problem schema to include %q", field)
		}
	}

	if _, ok := problemSchema.Properties["details"]; ok {
		t.Fatal("expected Problem schema to omit details")
	}
}

func newTestApp(t *testing.T, deps Dependencies, register ...func(huma.API)) *App {
	t.Helper()

	app, err := New(deps)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	for _, registerAPI := range register {
		registerAPI(app.API)
	}

	return app
}

func performRequest(app *App, method string, target string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, nil)
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, request)
	return recorder
}

func decodeProblem(t *testing.T, recorder *httptest.ResponseRecorder) Problem {
	t.Helper()

	var problem Problem
	if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
		t.Fatalf("decode problem response: %v", err)
	}

	return problem
}
