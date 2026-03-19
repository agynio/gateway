package handlers

import (
	"net/http"

	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/ziticonn"
	"github.com/agynio/gateway/internal/zitimgmtclient"
)

func NewAuthMiddleware(resolver zitimgmtclient.Resolver) func(http.Handler) http.Handler {
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
				problem := NewProblem(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "identity not available")
				WriteProblem(w, problem)
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
