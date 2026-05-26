package gateway

import (
	"context"

	"github.com/agynio/gateway/internal/identity"
)

func downstreamContext(ctx context.Context) context.Context {
	return identity.AppendToOutgoingContext(ctx)
}
