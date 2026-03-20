package gateway

import (
	"context"
	"errors"
	"io"

	"connectrpc.com/connect"
	notificationsv1 "github.com/agynio/gateway/gen/agynio/api/notifications/v1"
)

func (g *Gateway) Subscribe(ctx context.Context, req *connect.Request[notificationsv1.SubscribeRequest], stream *connect.ServerStream[notificationsv1.SubscribeResponse]) error {
	grpcStream, err := g.notifications.Subscribe(ctx, req.Msg)
	if err != nil {
		return toConnectError(err)
	}

	for {
		msg, err := grpcStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return toConnectError(err)
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
}
