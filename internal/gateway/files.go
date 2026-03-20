package gateway

import (
	"context"

	"connectrpc.com/connect"
	filesv1 "github.com/agynio/gateway/gen/agynio/api/files/v1"
)

func (g *Gateway) UploadFile(ctx context.Context, stream *connect.ClientStream[filesv1.UploadFileRequest]) (*connect.Response[filesv1.UploadFileResponse], error) {
	grpcStream, err := g.files.UploadFile(ctx)
	if err != nil {
		return nil, toConnectError(err)
	}

	for stream.Receive() {
		if err := grpcStream.Send(stream.Msg()); err != nil {
			return nil, toConnectError(err)
		}
	}
	if err := stream.Err(); err != nil {
		return nil, toConnectError(err)
	}

	resp, err := grpcStream.CloseAndRecv()
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetFileMetadata(ctx context.Context, req *connect.Request[filesv1.GetFileMetadataRequest]) (*connect.Response[filesv1.GetFileMetadataResponse], error) {
	resp, err := g.files.GetFileMetadata(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetDownloadUrl(ctx context.Context, req *connect.Request[filesv1.GetDownloadUrlRequest]) (*connect.Response[filesv1.GetDownloadUrlResponse], error) {
	resp, err := g.files.GetDownloadUrl(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
