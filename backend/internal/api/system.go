package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type healthStatus struct {
	Status string `json:"status" example:"ok"`
}

type healthResponse struct {
	Body healthStatus
}

func registerSystemRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "tonaris.system.health",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Report service health",
	}, func(context.Context, *struct{}) (*healthResponse, error) {
		return &healthResponse{
			Body: healthStatus{Status: "ok"},
		}, nil
	})
}
