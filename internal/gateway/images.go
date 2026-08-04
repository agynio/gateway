package gateway

import (
	"context"

	"connectrpc.com/connect"
	imagesv1 "github.com/agynio/gateway/gen/agynio/api/images/v1"
)

// ResolveVersion and RegisterPlatformImage are deliberately absent: both are
// internal, reached over the mesh rather than through the Gateway.

func (g *Gateway) CreateImage(ctx context.Context, req *connect.Request[imagesv1.CreateImageRequest]) (*connect.Response[imagesv1.CreateImageResponse], error) {
	resp, err := g.images.CreateImage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetImage(ctx context.Context, req *connect.Request[imagesv1.GetImageRequest]) (*connect.Response[imagesv1.GetImageResponse], error) {
	resp, err := g.images.GetImage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateImage(ctx context.Context, req *connect.Request[imagesv1.UpdateImageRequest]) (*connect.Response[imagesv1.UpdateImageResponse], error) {
	resp, err := g.images.UpdateImage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteImage(ctx context.Context, req *connect.Request[imagesv1.DeleteImageRequest]) (*connect.Response[imagesv1.DeleteImageResponse], error) {
	resp, err := g.images.DeleteImage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListImages(ctx context.Context, req *connect.Request[imagesv1.ListImagesRequest]) (*connect.Response[imagesv1.ListImagesResponse], error) {
	resp, err := g.images.ListImages(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListVersions(ctx context.Context, req *connect.Request[imagesv1.ListVersionsRequest]) (*connect.Response[imagesv1.ListVersionsResponse], error) {
	resp, err := g.images.ListVersions(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) RefreshImage(ctx context.Context, req *connect.Request[imagesv1.RefreshImageRequest]) (*connect.Response[imagesv1.RefreshImageResponse], error) {
	resp, err := g.images.RefreshImage(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
