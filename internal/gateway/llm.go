package gateway

import (
	"context"
	"errors"
	"io"

	"connectrpc.com/connect"
	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
)

func (g *Gateway) CreateResponse(ctx context.Context, req *connect.Request[llmv1.CreateResponseRequest]) (*connect.Response[llmv1.CreateResponseResponse], error) {
	resp, err := g.llm.CreateResponse(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateResponseStream(ctx context.Context, req *connect.Request[llmv1.CreateResponseStreamRequest], stream *connect.ServerStream[llmv1.CreateResponseStreamResponse]) error {
	grpcStream, err := g.llm.CreateResponseStream(ctx, req.Msg)
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
