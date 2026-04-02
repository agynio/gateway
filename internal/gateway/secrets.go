package gateway

import (
	"context"

	"connectrpc.com/connect"
	secretsv1 "github.com/agynio/gateway/gen/agynio/api/secrets/v1"
)

func (g *Gateway) ResolveSecret(ctx context.Context, req *connect.Request[secretsv1.ResolveSecretRequest]) (*connect.Response[secretsv1.ResolveSecretResponse], error) {
	resp, err := g.secrets.ResolveSecret(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateSecretProvider(ctx context.Context, req *connect.Request[secretsv1.CreateSecretProviderRequest]) (*connect.Response[secretsv1.CreateSecretProviderResponse], error) {
	resp, err := g.secrets.CreateSecretProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetSecretProvider(ctx context.Context, req *connect.Request[secretsv1.GetSecretProviderRequest]) (*connect.Response[secretsv1.GetSecretProviderResponse], error) {
	resp, err := g.secrets.GetSecretProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateSecretProvider(ctx context.Context, req *connect.Request[secretsv1.UpdateSecretProviderRequest]) (*connect.Response[secretsv1.UpdateSecretProviderResponse], error) {
	resp, err := g.secrets.UpdateSecretProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteSecretProvider(ctx context.Context, req *connect.Request[secretsv1.DeleteSecretProviderRequest]) (*connect.Response[secretsv1.DeleteSecretProviderResponse], error) {
	resp, err := g.secrets.DeleteSecretProvider(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSecretProviders(ctx context.Context, req *connect.Request[secretsv1.ListSecretProvidersRequest]) (*connect.Response[secretsv1.ListSecretProvidersResponse], error) {
	resp, err := g.secrets.ListSecretProviders(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) CreateSecret(ctx context.Context, req *connect.Request[secretsv1.CreateSecretRequest]) (*connect.Response[secretsv1.CreateSecretResponse], error) {
	resp, err := g.secrets.CreateSecret(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) GetSecret(ctx context.Context, req *connect.Request[secretsv1.GetSecretRequest]) (*connect.Response[secretsv1.GetSecretResponse], error) {
	resp, err := g.secrets.GetSecret(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) UpdateSecret(ctx context.Context, req *connect.Request[secretsv1.UpdateSecretRequest]) (*connect.Response[secretsv1.UpdateSecretResponse], error) {
	resp, err := g.secrets.UpdateSecret(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) DeleteSecret(ctx context.Context, req *connect.Request[secretsv1.DeleteSecretRequest]) (*connect.Response[secretsv1.DeleteSecretResponse], error) {
	resp, err := g.secrets.DeleteSecret(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func (g *Gateway) ListSecrets(ctx context.Context, req *connect.Request[secretsv1.ListSecretsRequest]) (*connect.Response[secretsv1.ListSecretsResponse], error) {
	resp, err := g.secrets.ListSecrets(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}
