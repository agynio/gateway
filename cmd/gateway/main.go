package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	agentstatev1 "github.com/agynio/gateway/gen/agynio/api/agent_state/v1"
	agentsv1 "github.com/agynio/gateway/gen/agynio/api/agents/v1"
	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
	filesv1 "github.com/agynio/gateway/gen/agynio/api/files/v1"
	"github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	notificationsv1 "github.com/agynio/gateway/gen/agynio/api/notifications/v1"
	secretsv1 "github.com/agynio/gateway/gen/agynio/api/secrets/v1"
	threadsv1 "github.com/agynio/gateway/gen/agynio/api/threads/v1"
	tokencountingv1 "github.com/agynio/gateway/gen/agynio/api/token_counting/v1"
	tracingv1 "github.com/agynio/gateway/gen/agynio/api/tracing/v1"
	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	zitimgmtv1 "github.com/agynio/gateway/gen/agynio/api/ziti_management/v1"
	"github.com/agynio/gateway/internal/apitokenresolver"
	"github.com/agynio/gateway/internal/clusteradminresolver"
	"github.com/agynio/gateway/internal/gateway"
	"github.com/agynio/gateway/internal/grpcclient"
	"github.com/agynio/gateway/internal/oidcauth"
	"github.com/agynio/gateway/internal/oidcresolver"
	"github.com/agynio/gateway/internal/platform"
	"github.com/agynio/gateway/internal/ziticonn"
	"github.com/agynio/gateway/internal/zitimgmtclient"
)

const (
	defaultAddr            = ":8080"
	defaultZitiServiceName = "gateway"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := platform.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load platform configuration: %v", err)
	}

	var zitiMgmtClient *zitimgmtclient.Client
	if config.ZitiEnabled {
		zitiMgmtClient, err = zitimgmtclient.NewClient(config.ZitiManagementGRPCTarget)
		if err != nil {
			log.Fatalf("failed to create ziti management gRPC client: %v", err)
		}
		defer func() {
			if err := zitiMgmtClient.Close(); err != nil {
				log.Printf("failed to close ziti management gRPC client: %v", err)
			}
		}()
	}

	var zitiResolver gateway.IdentityResolver
	if zitiMgmtClient != nil {
		zitiResolver = zitiMgmtClient
	}

	cleanup := make([]func(), 0, 10)
	defer func() {
		for _, closeFn := range cleanup {
			closeFn()
		}
	}()

	usersClient := mustClient(config.UsersGRPCTarget, "users", usersv1.NewUsersServiceClient, &cleanup)

	var oidcResolver gateway.BearerTokenResolver
	oidcConfigured := config.OIDCIssuerURL != "" || config.OIDCClientID != ""
	if oidcConfigured {
		if config.OIDCIssuerURL == "" || config.OIDCClientID == "" {
			log.Fatalf("both OIDC issuer URL and client ID are required when OIDC is enabled")
		}
		verifier, err := oidcauth.NewVerifier(ctx, config.OIDCIssuerURL, config.OIDCClientID)
		if err != nil {
			log.Fatalf("failed to create OIDC verifier: %v", err)
		}
		userinfoHTTPClient := &http.Client{Timeout: 10 * time.Second}
		resolver, err := oidcresolver.NewResolver(verifier, usersClient, userinfoHTTPClient)
		if err != nil {
			log.Fatalf("failed to create OIDC resolver: %v", err)
		}
		oidcResolver = resolver
	}

	apiTokenResolver := apitokenresolver.NewResolver(usersClient)

	var clusterAdminResolver *clusteradminresolver.Resolver
	if config.ClusterAdminToken != "" {
		resolver, err := clusteradminresolver.NewResolver(config.ClusterAdminToken, config.ClusterAdminIdentityID)
		if err != nil {
			log.Fatalf("failed to create cluster admin resolver: %v", err)
		}
		clusterAdminResolver = resolver
	}

	agentsClient := mustClient(config.AgentsGRPCTarget, "agents", agentsv1.NewAgentsServiceClient, &cleanup)
	threadsClient := mustClient(config.ThreadsGRPCTarget, "threads", threadsv1.NewThreadsServiceClient, &cleanup)
	chatClient := mustClient(config.ChatGRPCTarget, "chat", chatv1.NewChatServiceClient, &cleanup)
	notificationsClient := mustClient(config.NotificationsGRPCTarget, "notifications", notificationsv1.NewNotificationsServiceClient, &cleanup)
	filesClient := mustClient(config.FilesGRPCTarget, "files", filesv1.NewFilesServiceClient, &cleanup)
	agentStateClient := mustClient(config.AgentStateGRPCTarget, "agent state", agentstatev1.NewAgentStateServiceClient, &cleanup)
	tokenCountingClient := mustClient(config.TokenCountingGRPCTarget, "token counting", tokencountingv1.NewTokenCountingServiceClient, &cleanup)
	secretsClient := mustClient(config.SecretsGRPCTarget, "secrets", secretsv1.NewSecretsServiceClient, &cleanup)
	tracingClient := mustClient(config.TracingGRPCTarget, "tracing", tracingv1.NewTracingServiceClient, &cleanup)

	gatewayHandler := gateway.New(
		agentsClient,
		threadsClient,
		chatClient,
		notificationsClient,
		filesClient,
		agentStateClient,
		tokenCountingClient,
		secretsClient,
		tracingClient,
	)
	threadsGateway := gateway.NewThreadsGateway(gatewayHandler)
	usersGateway := gateway.NewUsersGateway(usersClient)

	interceptors := []connect.Interceptor{
		gateway.NewRecoveryInterceptor(),
		gateway.NewLoggingInterceptor(),
	}
	if zitiResolver != nil || oidcResolver != nil || apiTokenResolver != nil || clusterAdminResolver != nil {
		interceptors = append([]connect.Interceptor{gateway.NewAuthInterceptor(zitiResolver, oidcResolver, apiTokenResolver, clusterAdminResolver)}, interceptors...)
	}
	handlerOptions := connect.WithInterceptors(interceptors...)

	mux := http.NewServeMux()

	var meHandler http.Handler = http.HandlerFunc(gateway.MeHandler)
	if zitiResolver != nil || oidcResolver != nil || apiTokenResolver != nil || clusterAdminResolver != nil {
		meHandler = gateway.NewAuthMiddleware(zitiResolver, oidcResolver, apiTokenResolver, clusterAdminResolver)(meHandler)
	}
	mux.Handle("/me", meHandler)

	registerConnect := func(handlerPath string, handler http.Handler) {
		mux.Handle(handlerPath, handler)
	}

	registerConnect(gatewayv1connect.NewAgentsGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewThreadsGatewayHandler(threadsGateway, handlerOptions))
	registerConnect(gatewayv1connect.NewChatGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewNotificationsGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewFilesGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewAgentStateGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewTokenCountingGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewSecretsGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewTracingGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewUsersGatewayHandler(usersGateway, handlerOptions))

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"*"},
	})

	addr := defaultAddr
	if v := strings.TrimSpace(os.Getenv("ADDR")); v != "" {
		addr = v
	}

	connContext := func(ctx context.Context, conn net.Conn) context.Context {
		sourceIdentity, ok := ziticonn.SourceIdentityFromConn(conn)
		if !ok {
			return ctx
		}
		return ziticonn.WithSourceIdentity(ctx, sourceIdentity)
	}

	h2cHandler := h2c.NewHandler(corsMiddleware.Handler(mux), &http2.Server{})
	server := &http.Server{
		Addr:        addr,
		Handler:     h2cHandler,
		ConnContext: connContext,
	}

	if config.ZitiEnabled {
		zitiIdentityID, identityJSON, err := zitiMgmtClient.RequestServiceIdentity(ctx, zitimgmtv1.ServiceType_SERVICE_TYPE_GATEWAY)
		if err != nil {
			log.Fatalf("failed to request ziti service identity: %v", err)
		}

		zitiConfig := &ziti.Config{}
		if err := json.Unmarshal(identityJSON, zitiConfig); err != nil {
			log.Fatalf("failed to parse ziti identity: %v", err)
		}

		zitiContext, err := ziti.NewContext(zitiConfig)
		if err != nil {
			log.Fatalf("failed to create ziti context: %v", err)
		}
		defer zitiContext.Close()

		go renewLease(ctx, zitiMgmtClient, zitiIdentityID, config.ZitiLeaseRenewalInterval)

		serviceName := zitiServiceName()
		listener, err := zitiContext.ListenWithOptions(serviceName, ziti.DefaultListenOptions())
		if err != nil {
			log.Fatalf("failed to listen on ziti service %s: %v", serviceName, err)
		}

		log.Printf("gateway listening on ziti service %s", serviceName)
		go func() {
			if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("ziti server stopped: %v", err)
			}
		}()
	}

	log.Printf("gateway listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}

func renewLease(ctx context.Context, client *zitimgmtclient.Client, identityID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}
			if err := client.ExtendIdentityLease(ctx, identityID); err != nil {
				log.Printf("failed to extend ziti lease: %v", err)
			}
		}
	}
}

func zitiServiceName() string {
	if v := strings.TrimSpace(os.Getenv("ZITI_SERVICE_NAME")); v != "" {
		return v
	}
	return defaultZitiServiceName
}

func mustClient[T any](target, name string, factory func(grpc.ClientConnInterface) T, cleanup *[]func()) T {
	client, err := grpcclient.New(target, factory)
	if err != nil {
		log.Fatalf("failed to create %s gRPC client: %v", name, err)
	}

	if cleanup != nil {
		*cleanup = append(*cleanup, func() {
			if err := client.Close(); err != nil {
				log.Printf("failed to close %s gRPC client: %v", name, err)
			}
		})
	}

	return client.Service()
}
