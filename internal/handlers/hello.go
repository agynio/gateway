package handlers

import (
	"context"
	"fmt"

	"github.com/agynio/gateway/internal/gen"
)

const (
	healthStatusOK = "ok"
)

type Hello struct{}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) HealthCheck(_ context.Context, _ gen.HealthCheckRequestObject) (gen.HealthCheckResponseObject, error) {
	return gen.HealthCheck200JSONResponse(gen.HealthStatus{Status: healthStatusOK}), nil
}

func (h *Hello) SayHello(_ context.Context, request gen.SayHelloRequestObject) (gen.SayHelloResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	message := fmt.Sprintf("Hello, %s", request.Body.Name)
	return gen.SayHello200JSONResponse(gen.HelloResponse{Message: message}), nil
}
