package handlers

import (
	"context"
	"net/http"

	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
)

type IdentityResolver interface {
	ResolveIdentity(ctx context.Context, sourceIdentity string) (identity.ResolvedIdentity, error)
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

			sourceIdentity, ok := ziticonn.SourceIdentityFromContext(r.Context())
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			resolved, err := resolver.ResolveIdentity(r.Context(), sourceIdentity)
			if err != nil {
				problemErr := grpcErrorToProblem(err)
				WriteProblem(w, problemErr.Problem)
				return
			}

			ctx := identity.WithIdentity(r.Context(), resolved)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
