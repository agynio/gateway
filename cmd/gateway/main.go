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
	"github.com/agynio/gateway/internal/platform"
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

	config, err := platform.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load platform configuration: %v", err)
	}

	client, err := platform.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create platform client: %v", err)
	}

	root := chi.NewRouter()
	root.Use(chimw.RequestID)
	root.Use(chimw.RealIP)
	root.Use(chimw.Recoverer)
	root.Use(chimw.Logger)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"*"},
	})
	root.Use(corsMiddleware.Handler)

	teamRouter := chi.NewRouter()

	requestValidator, err := handlers.NewRequestValidationMiddleware(swagger)
	if err != nil {
		log.Fatalf("failed to initialise request validation: %v", err)
	}
	teamRouter.Use(requestValidator)

	if isResponseValidationEnabled() {
		responseValidator, err := handlers.NewResponseValidationMiddleware(swagger)
		if err != nil {
			log.Fatalf("failed to initialise response validation: %v", err)
		}
		teamRouter.Use(responseValidator)
	}

	strictHandler := gen.NewStrictHandlerWithOptions(handlers.NewTeam(client), nil, gen.StrictHTTPServerOptions{
		ResponseErrorHandlerFunc: handlers.StrictErrorHandler,
	})
	gen.HandlerWithOptions(strictHandler, gen.ChiServerOptions{BaseRouter: teamRouter})

	root.Mount(handlers.TeamBasePath(), teamRouter)

	addr := defaultAddr
	if v := strings.TrimSpace(os.Getenv("ADDR")); v != "" {
		addr = v
	}

	log.Printf("gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, root); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	swagger.Servers = []*openapi3.Server{{URL: handlers.TeamBasePath()}}
	return swagger, nil
}

func isResponseValidationEnabled() bool {
	return strings.EqualFold(os.Getenv("OPENAPI_VALIDATE_RESPONSE"), "true")
}
