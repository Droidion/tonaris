package api

import (
	"net/http"

	"backend/internal/authn"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

const (
	apiTitle   = "Tonaris API"
	apiVersion = "0.0.1"
)

type Dependencies struct {
	AllowedOrigins []string
	AuthVerifier   authn.Verifier
}

type App struct {
	Router http.Handler
	API    huma.API
}

func New(deps Dependencies) (*App, error) {
	configureErrorResponses()

	router := chi.NewMux()
	applyMiddleware(router, deps)

	config := huma.DefaultConfig(apiTitle, apiVersion)
	config.OpenAPIPath = "/api-doc/openapi"
	config.DocsPath = "/scalar"

	api := humachi.New(router, config)
	registerHelloRoutes(api)
	registerSystemRoutes(api)

	return &App{
		Router: router,
		API:    api,
	}, nil
}
