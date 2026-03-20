package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/agynio/gateway/gen/agynio/api/gateway/v1/gatewayv1connect"
	"github.com/agynio/gateway/internal/agentsclient"
	"github.com/agynio/gateway/internal/agentstateclient"
	"github.com/agynio/gateway/internal/chatclient"
	"github.com/agynio/gateway/internal/filesclient"
	"github.com/agynio/gateway/internal/gateway"
	"github.com/agynio/gateway/internal/llmclient"
	"github.com/agynio/gateway/internal/notificationsclient"
	"github.com/agynio/gateway/internal/platform"
	"github.com/agynio/gateway/internal/secretsclient"
	"github.com/agynio/gateway/internal/threadsclient"
	"github.com/agynio/gateway/internal/tokencountingclient"
	"github.com/agynio/gateway/internal/ziticonn"
	"github.com/agynio/gateway/internal/zitimgmtclient"
)

const (
	defaultAddr            = ":8080"
	defaultZitiServiceName = "gateway"
)

func main() {
	config, err := platform.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load platform configuration: %v", err)
	}

	var zitiMgmtClient *zitimgmtclient.Client
	if config.ZitiIdentityFile != "" {
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

	agentsClient, err := agentsclient.NewClient(config.AgentsGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create agents gRPC client: %v", err)
	}
	defer func() {
		if err := agentsClient.Close(); err != nil {
			log.Printf("failed to close agents gRPC client: %v", err)
		}
	}()

	threadsClient, err := threadsclient.NewClient(config.ThreadsGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create threads gRPC client: %v", err)
	}
	defer func() {
		if err := threadsClient.Close(); err != nil {
			log.Printf("failed to close threads gRPC client: %v", err)
		}
	}()

	chatClient, err := chatclient.NewClient(config.ChatGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create chat gRPC client: %v", err)
	}
	defer func() {
		if err := chatClient.Close(); err != nil {
			log.Printf("failed to close chat gRPC client: %v", err)
		}
	}()

	notificationsClient, err := notificationsclient.NewClient(config.NotificationsGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create notifications gRPC client: %v", err)
	}
	defer func() {
		if err := notificationsClient.Close(); err != nil {
			log.Printf("failed to close notifications gRPC client: %v", err)
		}
	}()

	filesClient, err := filesclient.NewClient(config.FilesGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create files gRPC client: %v", err)
	}
	defer func() {
		if err := filesClient.Close(); err != nil {
			log.Printf("failed to close files gRPC client: %v", err)
		}
	}()

	agentStateClient, err := agentstateclient.NewClient(config.AgentStateGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create agent state gRPC client: %v", err)
	}
	defer func() {
		if err := agentStateClient.Close(); err != nil {
			log.Printf("failed to close agent state gRPC client: %v", err)
		}
	}()

	tokenCountingClient, err := tokencountingclient.NewClient(config.TokenCountingGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create token counting gRPC client: %v", err)
	}
	defer func() {
		if err := tokenCountingClient.Close(); err != nil {
			log.Printf("failed to close token counting gRPC client: %v", err)
		}
	}()

	llmClient, err := llmclient.NewClient(config.LLMGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create llm gRPC client: %v", err)
	}
	defer func() {
		if err := llmClient.Close(); err != nil {
			log.Printf("failed to close llm gRPC client: %v", err)
		}
	}()

	secretsClient, err := secretsclient.NewClient(config.SecretsGRPCTarget)
	if err != nil {
		log.Fatalf("failed to create secrets gRPC client: %v", err)
	}
	defer func() {
		if err := secretsClient.Close(); err != nil {
			log.Printf("failed to close secrets gRPC client: %v", err)
		}
	}()

	gatewayHandler := gateway.New(
		agentsClient.AgentsServiceClient(),
		threadsClient.ThreadsServiceClient(),
		notificationsClient.NotificationsServiceClient(),
		filesClient.FilesServiceClient(),
		agentStateClient.AgentStateServiceClient(),
		tokenCountingClient.TokenCountingServiceClient(),
		llmClient.LLMServiceClient(),
		secretsClient.SecretsServiceClient(),
	)
	chatGateway := gateway.NewChatGateway(chatClient.ChatServiceClient())

	interceptors := []connect.Interceptor{
		gateway.NewRecoveryInterceptor(),
		gateway.NewLoggingInterceptor(),
	}
	if zitiMgmtClient != nil {
		interceptors = append([]connect.Interceptor{gateway.NewAuthInterceptor(zitiMgmtClient)}, interceptors...)
	}
	handlerOptions := connect.WithInterceptors(interceptors...)

	mux := http.NewServeMux()

	var meHandler http.Handler = http.HandlerFunc(gateway.MeHandler)
	if zitiMgmtClient != nil {
		meHandler = gateway.NewAuthMiddleware(zitiMgmtClient)(meHandler)
	}
	mux.Handle("/me", meHandler)

	registerConnect := func(handlerPath string, handler http.Handler) {
		mux.Handle(handlerPath, handler)
	}

	registerConnect(gatewayv1connect.NewAgentsGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewThreadsGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewChatGatewayHandler(chatGateway, handlerOptions))
	registerConnect(gatewayv1connect.NewNotificationsGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewFilesGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewAgentStateGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewTokenCountingGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewLLMGatewayHandler(gatewayHandler, handlerOptions))
	registerConnect(gatewayv1connect.NewSecretsGatewayHandler(gatewayHandler, handlerOptions))

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

	if config.ZitiIdentityFile != "" {
		zitiContext, err := ziti.NewContextFromFile(config.ZitiIdentityFile)
		if err != nil {
			log.Fatalf("failed to create ziti context: %v", err)
		}
		defer zitiContext.Close()

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

func zitiServiceName() string {
	if v := strings.TrimSpace(os.Getenv("ZITI_SERVICE_NAME")); v != "" {
		return v
	}
	return defaultZitiServiceName
}
