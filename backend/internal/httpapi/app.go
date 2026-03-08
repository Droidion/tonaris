package httpapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

const (
	apiTitle   = "Tonaris API"
	apiVersion = "0.0.1"
)

type App struct {
	Router http.Handler
	API    huma.API
}

func New() (*App, error) {
	configureErrorResponses()

	router := chi.NewMux()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	config := huma.DefaultConfig(apiTitle, apiVersion)
	config.OpenAPIPath = "/api-doc/openapi"
	config.DocsPath = "/scalar"

	api := humachi.New(router, config)

	huma.Register(api, huma.Operation{
		OperationID: "tonaris.hello",
		Method:      http.MethodGet,
		Path:        "/hello",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Ok",
				Content: map[string]*huma.MediaType{
					"text/plain": {
						Schema: &huma.Schema{Type: "string"},
					},
				},
			},
		},
	}, func(context.Context, *struct{}) (*helloResponse, error) {
		return &helloResponse{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		}, nil
	})

	return &App{
		Router: router,
		API:    api,
	}, nil
}

type helloResponse struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}
