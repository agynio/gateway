package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"connectrpc.com/connect"
	"github.com/agynio/gateway/internal/apitokenresolver"
	"github.com/agynio/gateway/internal/httpauth"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityResolver interface {
	ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error)
}

type BearerTokenResolver interface {
	ResolveFromToken(ctx context.Context, accessToken string) (identity.ResolvedIdentity, error)
}

func NewAuthInterceptor(zitiResolver IdentityResolver, oidcResolver BearerTokenResolver, apiTokenResolver BearerTokenResolver) connect.Interceptor {
	if zitiResolver == nil && oidcResolver == nil && apiTokenResolver == nil {
		panic("at least one identity resolver is required")
	}
	return authInterceptor{zitiResolver: zitiResolver, oidcResolver: oidcResolver, apiTokenResolver: apiTokenResolver}
}

func NewAuthMiddleware(zitiResolver IdentityResolver, oidcResolver BearerTokenResolver, apiTokenResolver BearerTokenResolver) func(http.Handler) http.Handler {
	if zitiResolver == nil && oidcResolver == nil && apiTokenResolver == nil {
		panic("at least one identity resolver is required")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			ctx := withBearerToken(r.Context(), r.Header.Get("Authorization"))
			ctx, err := resolveIdentity(ctx, zitiResolver, oidcResolver, apiTokenResolver)
			if err != nil {
				writeProblem(w, httpStatusFromError(err), err.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type authInterceptor struct {
	zitiResolver     IdentityResolver
	oidcResolver     BearerTokenResolver
	apiTokenResolver BearerTokenResolver
}

func (a authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx = withBearerToken(ctx, req.Header().Get("Authorization"))
		ctx, err := resolveIdentity(ctx, a.zitiResolver, a.oidcResolver, a.apiTokenResolver)
		if err != nil {
			return nil, toConnectError(err)
		}
		return next(ctx, req)
	}
}

func (a authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		ctx = withBearerToken(ctx, conn.RequestHeader().Get("Authorization"))
		ctx, err := resolveIdentity(ctx, a.zitiResolver, a.oidcResolver, a.apiTokenResolver)
		if err != nil {
			return toConnectError(err)
		}
		return next(ctx, conn)
	}
}

type recoveryInterceptor struct{}

func NewRecoveryInterceptor() connect.Interceptor {
	return recoveryInterceptor{}
}

func (recoveryInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic in %s: %v\n%s", req.Spec().Procedure, recovered, debug.Stack())
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("internal error"))
			}
		}()
		return next(ctx, req)
	}
}

func (recoveryInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (recoveryInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic in %s: %v\n%s", conn.Spec().Procedure, recovered, debug.Stack())
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("internal error"))
			}
		}()
		return next(ctx, conn)
	}
}

type loggingInterceptor struct{}

func NewLoggingInterceptor() connect.Interceptor {
	return loggingInterceptor{}
}

func (loggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		resp, err := next(ctx, req)
		duration := time.Since(start)
		if err != nil {
			log.Printf("rpc %s failed in %s: %v", req.Spec().Procedure, duration, err)
			return resp, err
		}
		log.Printf("rpc %s completed in %s", req.Spec().Procedure, duration)
		return resp, nil
	}
}

func (loggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (loggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		start := time.Now()
		err := next(ctx, conn)
		duration := time.Since(start)
		if err != nil {
			log.Printf("rpc %s failed in %s: %v", conn.Spec().Procedure, duration, err)
			return err
		}
		log.Printf("rpc %s completed in %s", conn.Spec().Procedure, duration)
		return nil
	}
}

func resolveIdentity(ctx context.Context, zitiResolver IdentityResolver, oidcResolver BearerTokenResolver, apiTokenResolver BearerTokenResolver) (context.Context, error) {
	sourceIdentity, ok := ziticonn.SourceIdentityFromContext(ctx)
	if ok {
		if zitiResolver == nil {
			return ctx, status.Error(codes.Unauthenticated, "ziti identity resolver is not configured")
		}
		resolvedIdentity, err := zitiResolver.ResolveIdentity(ctx, sourceIdentity)
		if err != nil {
			return ctx, err
		}
		return identity.WithIdentity(ctx, resolvedIdentity), nil
	}

	accessToken, ok := httpauth.BearerTokenFromContext(ctx)
	if ok {
		if apitokenresolver.HasPrefix(accessToken) {
			if apiTokenResolver == nil {
				return ctx, status.Error(codes.Unauthenticated, "api token resolver is not configured")
			}
			resolvedIdentity, err := apiTokenResolver.ResolveFromToken(ctx, accessToken)
			if err != nil {
				return ctx, err
			}
			return identity.WithIdentity(ctx, resolvedIdentity), nil
		}
		if oidcResolver == nil {
			return ctx, status.Error(codes.Unauthenticated, "oidc resolver is not configured")
		}
		resolvedIdentity, err := oidcResolver.ResolveFromToken(ctx, accessToken)
		if err != nil {
			return ctx, err
		}
		return identity.WithIdentity(ctx, resolvedIdentity), nil
	}

	return ctx, nil
}

func withBearerToken(ctx context.Context, authHeader string) context.Context {
	accessToken, ok := httpauth.ExtractBearerToken(authHeader)
	if !ok {
		return ctx
	}
	return httpauth.WithBearerToken(ctx, accessToken)
}
