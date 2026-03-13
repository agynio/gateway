// All API domains must be wired through oapi-codegen generated strict servers
// with request validation middleware. Do NOT register raw http.Handler routes
// for CRUD endpoints. The only exception is streaming/proxy endpoints (e.g.,
// SSE /responses) which must still be defined in the OpenAPI spec for request
// validation but mounted as raw handlers outside the strict server group.
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

	llmv1schema "github.com/agynio/gateway/internal/apischema/llmv1"
	teamv1schema "github.com/agynio/gateway/internal/apischema/teamv1"
	"github.com/agynio/gateway/internal/filesclient"
	"github.com/agynio/gateway/internal/gen"
	"github.com/agynio/gateway/internal/handlers"
	"github.com/agynio/gateway/internal/llmclient"
	"github.com/agynio/gateway/internal/llmgen"
	"github.com/agynio/gateway/internal/platform"
	"github.com/agynio/gateway/internal/secretsclient"
	"github.com/agynio/gateway/internal/teamsclient"
)

const (
	defaultAddr = ":8080"
)

func main() {
	teamSpec, err := loadTeamSpec()
	if err != nil {
		log.Fatalf("failed to load team OpenAPI spec: %v", err)
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

	requestValidator, err := handlers.NewRequestValidationMiddleware(teamSpec)
	if err != nil {
		log.Fatalf("failed to initialise request validation: %v", err)
	}
	teamRouter.Use(requestValidator)

	if isResponseValidationEnabled() {
		responseValidator, err := handlers.NewResponseValidationMiddleware(teamSpec)
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

	llmGRPCEnabled := config.LLMGRPCTarget != ""
	llmHTTPEnabled := config.LLMHTTPBaseURL != nil

	if llmGRPCEnabled || llmHTTPEnabled {
		llmSpec, err := loadLLMSpec()
		if err != nil {
			log.Fatalf("failed to load llm OpenAPI spec: %v", err)
		}

		llmRequestValidator, err := handlers.NewRequestValidationMiddleware(llmSpec)
		if err != nil {
			log.Fatalf("failed to initialise llm request validation: %v", err)
		}

		llmRouter := chi.NewRouter()
		llmRouter.Use(llmRequestValidator)

		if llmGRPCEnabled {
			llmClient, err := llmclient.NewClient(config.LLMGRPCTarget)
			if err != nil {
				log.Fatalf("failed to create llm gRPC client: %v", err)
			}
			defer func() {
				if err := llmClient.Close(); err != nil {
					log.Printf("failed to close llm gRPC client: %v", err)
				}
			}()

			llmHandler := handlers.NewLLMHandler(llmClient)
			llmStrictHandler := llmgen.NewStrictHandlerWithOptions(llmHandler, nil, llmgen.StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  handlers.StrictRequestErrorHandler,
				ResponseErrorHandlerFunc: handlers.StrictErrorHandler,
			})

			var llmResponseValidator func(http.Handler) http.Handler
			if isResponseValidationEnabled() {
				llmResponseValidator, err = handlers.NewResponseValidationMiddleware(llmSpec)
				if err != nil {
					log.Fatalf("failed to initialise llm response validation: %v", err)
				}
			}

			llmRouter.Group(func(r chi.Router) {
				if llmResponseValidator != nil {
					r.Use(llmResponseValidator)
				}
				llmgen.HandlerWithOptions(llmStrictHandler, llmgen.ChiServerOptions{BaseRouter: r})
			})
		}

		if llmHTTPEnabled {
			llmProxy := handlers.NewLLMResponseProxy(config.LLMHTTPBaseURL)
			llmRouter.Post("/responses", llmProxy.ServeHTTP)
		}

		root.Mount(handlers.LLMBasePath(), llmRouter)
	}

	if config.SecretsGRPCTarget != "" {
		secretsClient, err := secretsclient.NewClient(config.SecretsGRPCTarget)
		if err != nil {
			log.Fatalf("failed to create secrets gRPC client: %v", err)
		}
		defer func() {
			if err := secretsClient.Close(); err != nil {
				log.Printf("failed to close secrets gRPC client: %v", err)
			}
		}()

		secretsHandler := handlers.NewSecretsHandler(secretsClient)
		root.Route("/secrets/v1", func(r chi.Router) {
			r.Post("/secret-providers", secretsHandler.CreateProvider)
			r.Get("/secret-providers", secretsHandler.ListProviders)
			r.Get("/secret-providers/{providerId}", secretsHandler.GetProvider)
			r.Patch("/secret-providers/{providerId}", secretsHandler.UpdateProvider)
			r.Delete("/secret-providers/{providerId}", secretsHandler.DeleteProvider)
			r.Post("/secrets", secretsHandler.CreateSecret)
			r.Get("/secrets", secretsHandler.ListSecrets)
			r.Get("/secrets/{secretId}", secretsHandler.GetSecret)
			r.Patch("/secrets/{secretId}", secretsHandler.UpdateSecret)
			r.Delete("/secrets/{secretId}", secretsHandler.DeleteSecret)
			r.Post("/secrets/{secretId}/resolve", secretsHandler.ResolveSecret)
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

func loadTeamSpec() (*openapi3.T, error) {
	swagger, err := teamv1schema.LoadSpec()
	if err != nil {
		return nil, err
	}
	swagger.Servers = []*openapi3.Server{{URL: handlers.TeamBasePath()}}
	return swagger, nil
}

func loadLLMSpec() (*openapi3.T, error) {
	swagger, err := llmv1schema.LoadSpec()
	if err != nil {
		return nil, err
	}
	swagger.Servers = []*openapi3.Server{{URL: handlers.LLMBasePath()}}
	return swagger, nil
}

func isResponseValidationEnabled() bool {
	return strings.EqualFold(os.Getenv("OPENAPI_VALIDATE_RESPONSE"), "true")
}
