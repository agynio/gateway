package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
)

type IdentityResolver interface {
	ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error)
}

func NewAuthInterceptor(resolver IdentityResolver) connect.Interceptor {
	if resolver == nil {
		panic("identity resolver is required")
	}
	return authInterceptor{resolver: resolver}
}

func NewAuthMiddleware(resolver IdentityResolver) func(http.Handler) http.Handler {
	if resolver == nil {
		panic("identity resolver is required")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			ctx, err := resolveIdentity(r.Context(), resolver)
			if err != nil {
				writeProblem(w, httpStatusFromError(err), err.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type authInterceptor struct {
	resolver IdentityResolver
}

func (a authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx, err := resolveIdentity(ctx, a.resolver)
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
		ctx, err := resolveIdentity(ctx, a.resolver)
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
				log.Printf("panic in %s: %v", req.Spec().Procedure, recovered)
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
				log.Printf("panic in %s: %v", conn.Spec().Procedure, recovered)
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

func resolveIdentity(ctx context.Context, resolver IdentityResolver) (context.Context, error) {
	sourceIdentity, ok := ziticonn.SourceIdentityFromContext(ctx)
	if !ok {
		return ctx, nil
	}

	resolvedIdentity, err := resolver.ResolveIdentity(ctx, sourceIdentity)
	if err != nil {
		return ctx, err
	}

	return identity.WithIdentity(ctx, resolvedIdentity), nil
}
