package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type helloResponse struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

func registerHelloRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "tonaris.system.hello",
		Method:      http.MethodGet,
		Path:        "/hello",
		Summary:     "Return a simple hello world response",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "OK",
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
			Body:        []byte("hello world"),
		}, nil
	})
}
