package gateway

import (
	"context"

	"connectrpc.com/connect"
	networksv1 "github.com/agynio/gateway/gen/agynio/api/networks/v1"
)

type NetworksGateway struct {
	networks networksv1.NetworksServiceClient
}

func NewNetworksGateway(networks networksv1.NetworksServiceClient) *NetworksGateway {
	return &NetworksGateway{networks: networks}
}

func (g *NetworksGateway) CreateNetwork(ctx context.Context, req *connect.Request[networksv1.CreateNetworkRequest]) (*connect.Response[networksv1.CreateNetworkResponse], error) {
	resp, err := g.networks.CreateNetwork(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) GetNetwork(ctx context.Context, req *connect.Request[networksv1.GetNetworkRequest]) (*connect.Response[networksv1.GetNetworkResponse], error) {
	resp, err := g.networks.GetNetwork(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) ListNetworks(ctx context.Context, req *connect.Request[networksv1.ListNetworksRequest]) (*connect.Response[networksv1.ListNetworksResponse], error) {
	resp, err := g.networks.ListNetworks(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) UpdateNetwork(ctx context.Context, req *connect.Request[networksv1.UpdateNetworkRequest]) (*connect.Response[networksv1.UpdateNetworkResponse], error) {
	resp, err := g.networks.UpdateNetwork(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) DeleteNetwork(ctx context.Context, req *connect.Request[networksv1.DeleteNetworkRequest]) (*connect.Response[networksv1.DeleteNetworkResponse], error) {
	resp, err := g.networks.DeleteNetwork(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) CreateTunnelCredential(ctx context.Context, req *connect.Request[networksv1.CreateTunnelCredentialRequest]) (*connect.Response[networksv1.CreateTunnelCredentialResponse], error) {
	resp, err := g.networks.CreateTunnelCredential(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) GetTunnelCredential(ctx context.Context, req *connect.Request[networksv1.GetTunnelCredentialRequest]) (*connect.Response[networksv1.GetTunnelCredentialResponse], error) {
	resp, err := g.networks.GetTunnelCredential(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) ListTunnelCredentials(ctx context.Context, req *connect.Request[networksv1.ListTunnelCredentialsRequest]) (*connect.Response[networksv1.ListTunnelCredentialsResponse], error) {
	resp, err := g.networks.ListTunnelCredentials(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) DeleteTunnelCredential(ctx context.Context, req *connect.Request[networksv1.DeleteTunnelCredentialRequest]) (*connect.Response[networksv1.DeleteTunnelCredentialResponse], error) {
	resp, err := g.networks.DeleteTunnelCredential(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) CreatePrivateResource(ctx context.Context, req *connect.Request[networksv1.CreatePrivateResourceRequest]) (*connect.Response[networksv1.CreatePrivateResourceResponse], error) {
	resp, err := g.networks.CreatePrivateResource(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) GetPrivateResource(ctx context.Context, req *connect.Request[networksv1.GetPrivateResourceRequest]) (*connect.Response[networksv1.GetPrivateResourceResponse], error) {
	resp, err := g.networks.GetPrivateResource(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) ListPrivateResources(ctx context.Context, req *connect.Request[networksv1.ListPrivateResourcesRequest]) (*connect.Response[networksv1.ListPrivateResourcesResponse], error) {
	resp, err := g.networks.ListPrivateResources(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) UpdatePrivateResource(ctx context.Context, req *connect.Request[networksv1.UpdatePrivateResourceRequest]) (*connect.Response[networksv1.UpdatePrivateResourceResponse], error) {
	resp, err := g.networks.UpdatePrivateResource(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) DeletePrivateResource(ctx context.Context, req *connect.Request[networksv1.DeletePrivateResourceRequest]) (*connect.Response[networksv1.DeletePrivateResourceResponse], error) {
	resp, err := g.networks.DeletePrivateResource(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) CreatePrivateResourceAccess(ctx context.Context, req *connect.Request[networksv1.CreatePrivateResourceAccessRequest]) (*connect.Response[networksv1.CreatePrivateResourceAccessResponse], error) {
	resp, err := g.networks.CreatePrivateResourceAccess(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) DeletePrivateResourceAccess(ctx context.Context, req *connect.Request[networksv1.DeletePrivateResourceAccessRequest]) (*connect.Response[networksv1.DeletePrivateResourceAccessResponse], error) {
	resp, err := g.networks.DeletePrivateResourceAccess(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *NetworksGateway) ListPrivateResourceAccess(ctx context.Context, req *connect.Request[networksv1.ListPrivateResourceAccessRequest]) (*connect.Response[networksv1.ListPrivateResourceAccessResponse], error) {
	resp, err := g.networks.ListPrivateResourceAccess(downstreamContext(ctx), req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
