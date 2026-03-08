package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var defaultAllowedOrigins = []string{"http://localhost:3000"}

func applyMiddleware(router chi.Router, deps Dependencies) {
	// Zero-value dependencies are accepted so tests and tooling can build the app
	// without runtime config or external integrations.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOriginsOrDefault(deps.AllowedOrigins),
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
}

func allowedOriginsOrDefault(origins []string) []string {
	if len(origins) == 0 {
		return append([]string(nil), defaultAllowedOrigins...)
	}

	return append([]string(nil), origins...)
}
