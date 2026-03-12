package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	teamv1schema "github.com/agynio/gateway/internal/apischema/teamv1"
	"github.com/agynio/gateway/internal/filesclient"
	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/handlers"
	"github.com/agynio/gateway/internal/platform"
	"github.com/agynio/gateway/internal/teamsclient"
)

const (
	defaultAddr = ":8080"
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

	teamsClient, err := teamsclient.NewClient(config.TeamsGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create teams gRPC client: %v", err)
	}
	defer func() {
		if err := teamsClient.Close(); err != nil {
			log.Printf("failed to close teams gRPC client: %v", err)
		}
	}()

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

	strictHandler := gen.NewStrictHandlerWithOptions(handlers.NewTeam(teamsClient.TeamsServiceClient()), nil, gen.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  handlers.StrictRequestErrorHandler,
		ResponseErrorHandlerFunc: handlers.StrictErrorHandler,
	})
	gen.HandlerWithOptions(strictHandler, gen.ChiServerOptions{BaseRouter: teamRouter})

	root.Mount(handlers.TeamBasePath(), teamRouter)

	if config.FilesGRPCTarget != "" {
		filesClient, err := filesclient.NewClient(config.FilesGRPCTarget)
		if err != nil {
			log.Fatalf("failed to create files gRPC client: %v", err)
		}
		defer func() {
			if err := filesClient.Close(); err != nil {
				log.Printf("failed to close files gRPC client: %v", err)
			}
		}()
		filesHandler := handlers.NewFilesHandler(filesClient)
		root.Route("/files/v1", func(r chi.Router) {
			r.Post("/files", filesHandler.Upload)
		})
	}

	proxyHandler := handlers.NewUpstreamProxy(client)
	root.Handle("/health", proxyHandler)
	root.Route("/api", func(r chi.Router) {
		r.Handle("/", proxyHandler)
		r.Handle("/*", proxyHandler)
	})

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
	swagger, err := teamv1schema.LoadSpec()
	if err != nil {
		return nil, err
	}
	swagger.Servers = []*openapi3.Server{{URL: handlers.TeamBasePath()}}
	return swagger, nil
}

func isResponseValidationEnabled() bool {
	return strings.EqualFold(os.Getenv("OPENAPI_VALIDATE_RESPONSE"), "true")
}
