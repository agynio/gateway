package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/handlers"
)

const (
	defaultAddr      = ":8080"
	specRelativePath = "spec/openapi.yaml"
)

func main() {
	swagger, err := loadSpec()
	if err != nil {
		log.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	router := chi.NewRouter()
	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Recoverer)
	router.Use(chimw.Logger)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"*"},
	})
	router.Use(corsMiddleware.Handler)

	if isResponseValidationEnabled() {
		responseValidator, err := handlers.NewResponseValidationMiddleware(swagger)
		if err != nil {
			log.Fatalf("failed to initialise response validation: %v", err)
		}
		router.Use(responseValidator)
	}

	requestValidator, err := handlers.NewRequestValidationMiddleware(swagger)
	if err != nil {
		log.Fatalf("failed to initialise request validation: %v", err)
	}
	router.Use(requestValidator)

	strictHandler := handlers.NewHello()
	server := gen.NewStrictHandler(strictHandler, nil)
	gen.HandlerWithOptions(server, gen.ChiServerOptions{BaseRouter: router})

	addr := defaultAddr
	if v := strings.TrimSpace(os.Getenv("ADDR")); v != "" {
		addr = v
	}

	log.Printf("gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}

func loadSpec() (*openapi3.T, error) {
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	specPath := filepath.Clean(specRelativePath)
	swagger, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil
	return swagger, nil
}

func isResponseValidationEnabled() bool {
	return strings.EqualFold(os.Getenv("VALIDATE_RESPONSES"), "true")
}
