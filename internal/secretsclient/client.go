package secretsclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	secretsv1 "github.com/agynio/gateway/gen/agynio/api/secrets/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SecretProviderType string

const (
	SecretProviderTypeVault SecretProviderType = "vault"
)

type Client struct {
	conn   *grpc.ClientConn
	client secretsv1.SecretsServiceClient
}

type VaultConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

type SecretProviderConfig struct {
	Vault *VaultConfig `json:"vault,omitempty"`
}

type SecretProvider struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Type        SecretProviderType   `json:"type"`
	Config      SecretProviderConfig `json:"config"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   *time.Time           `json:"updatedAt,omitempty"`
}

type Secret struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	SecretProviderID string     `json:"secretProviderId"`
	RemoteName       string     `json:"remoteName"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

type ResolvedSecret struct {
	Value string `json:"value"`
}

type CreateProviderParams struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Type        SecretProviderType   `json:"type"`
	Config      SecretProviderConfig `json:"config"`
}

type UpdateProviderParams struct {
	Title       *string               `json:"title"`
	Description *string               `json:"description"`
	Config      *SecretProviderConfig `json:"config"`
}

type CreateSecretParams struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	SecretProviderID string `json:"secretProviderId"`
	RemoteName       string `json:"remoteName"`
}

type UpdateSecretParams struct {
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	SecretProviderID *string `json:"secretProviderId"`
	RemoteName       *string `json:"remoteName"`
}

func NewClient(target string) (*Client, error) {
	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("target is required")
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: secretsv1.NewSecretsServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateProvider(ctx context.Context, params CreateProviderParams) (SecretProvider, error) {
	providerType, err := toProtoProviderType(params.Type)
	if err != nil {
		return SecretProvider{}, err
	}

	config, err := toProtoProviderConfig(params.Config)
	if err != nil {
		return SecretProvider{}, err
	}

	resp, err := c.client.CreateSecretProvider(ctx, &secretsv1.CreateSecretProviderRequest{
		Title:       params.Title,
		Description: params.Description,
		Type:        providerType,
		Config:      config,
	})
	if err != nil {
		return SecretProvider{}, err
	}

	provider := resp.GetSecretProvider()
	if provider == nil {
		return SecretProvider{}, fmt.Errorf("create provider response missing provider")
	}

	return providerFromProto(provider)
}

func (c *Client) GetProvider(ctx context.Context, id string) (SecretProvider, error) {
	resp, err := c.client.GetSecretProvider(ctx, &secretsv1.GetSecretProviderRequest{Id: id})
	if err != nil {
		return SecretProvider{}, err
	}

	provider := resp.GetSecretProvider()
	if provider == nil {
		return SecretProvider{}, fmt.Errorf("get provider response missing provider")
	}

	return providerFromProto(provider)
}

func (c *Client) UpdateProvider(ctx context.Context, id string, params UpdateProviderParams) (SecretProvider, error) {
	req := &secretsv1.UpdateSecretProviderRequest{Id: id}
	if params.Title != nil {
		req.Title = params.Title
	}
	if params.Description != nil {
		req.Description = params.Description
	}
	if params.Config != nil {
		config, err := toProtoProviderConfig(*params.Config)
		if err != nil {
			return SecretProvider{}, err
		}
		req.Config = config
	}

	resp, err := c.client.UpdateSecretProvider(ctx, req)
	if err != nil {
		return SecretProvider{}, err
	}

	provider := resp.GetSecretProvider()
	if provider == nil {
		return SecretProvider{}, fmt.Errorf("update provider response missing provider")
	}

	return providerFromProto(provider)
}

func (c *Client) DeleteProvider(ctx context.Context, id string) error {
	_, err := c.client.DeleteSecretProvider(ctx, &secretsv1.DeleteSecretProviderRequest{Id: id})
	return err
}

func (c *Client) ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]SecretProvider, string, error) {
	resp, err := c.client.ListSecretProviders(ctx, &secretsv1.ListSecretProvidersRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	})
	if err != nil {
		return nil, "", err
	}

	providers := make([]SecretProvider, 0, len(resp.GetSecretProviders()))
	for _, provider := range resp.GetSecretProviders() {
		converted, err := providerFromProto(provider)
		if err != nil {
			return nil, "", err
		}
		providers = append(providers, converted)
	}
	if providers == nil {
		providers = []SecretProvider{}
	}

	return providers, resp.GetNextPageToken(), nil
}

func (c *Client) CreateSecret(ctx context.Context, params CreateSecretParams) (Secret, error) {
	resp, err := c.client.CreateSecret(ctx, &secretsv1.CreateSecretRequest{
		Title:            params.Title,
		Description:      params.Description,
		SecretProviderId: params.SecretProviderID,
		RemoteName:       params.RemoteName,
	})
	if err != nil {
		return Secret{}, err
	}

	secret := resp.GetSecret()
	if secret == nil {
		return Secret{}, fmt.Errorf("create secret response missing secret")
	}

	return secretFromProto(secret)
}

func (c *Client) GetSecret(ctx context.Context, id string) (Secret, error) {
	resp, err := c.client.GetSecret(ctx, &secretsv1.GetSecretRequest{Id: id})
	if err != nil {
		return Secret{}, err
	}

	secret := resp.GetSecret()
	if secret == nil {
		return Secret{}, fmt.Errorf("get secret response missing secret")
	}

	return secretFromProto(secret)
}

func (c *Client) UpdateSecret(ctx context.Context, id string, params UpdateSecretParams) (Secret, error) {
	req := &secretsv1.UpdateSecretRequest{Id: id}
	if params.Title != nil {
		req.Title = params.Title
	}
	if params.Description != nil {
		req.Description = params.Description
	}
	if params.SecretProviderID != nil {
		req.SecretProviderId = params.SecretProviderID
	}
	if params.RemoteName != nil {
		req.RemoteName = params.RemoteName
	}

	resp, err := c.client.UpdateSecret(ctx, req)
	if err != nil {
		return Secret{}, err
	}

	secret := resp.GetSecret()
	if secret == nil {
		return Secret{}, fmt.Errorf("update secret response missing secret")
	}

	return secretFromProto(secret)
}

func (c *Client) DeleteSecret(ctx context.Context, id string) error {
	_, err := c.client.DeleteSecret(ctx, &secretsv1.DeleteSecretRequest{Id: id})
	return err
}

func (c *Client) ListSecrets(ctx context.Context, pageSize int32, pageToken, secretProviderID string) ([]Secret, string, error) {
	req := &secretsv1.ListSecretsRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}
	if strings.TrimSpace(secretProviderID) != "" {
		req.SecretProviderId = secretProviderID
	}

	resp, err := c.client.ListSecrets(ctx, req)
	if err != nil {
		return nil, "", err
	}

	secrets := make([]Secret, 0, len(resp.GetSecrets()))
	for _, secret := range resp.GetSecrets() {
		converted, err := secretFromProto(secret)
		if err != nil {
			return nil, "", err
		}
		secrets = append(secrets, converted)
	}
	if secrets == nil {
		secrets = []Secret{}
	}

	return secrets, resp.GetNextPageToken(), nil
}

func (c *Client) ResolveSecret(ctx context.Context, id string) (ResolvedSecret, error) {
	resp, err := c.client.ResolveSecret(ctx, &secretsv1.ResolveSecretRequest{Id: id})
	if err != nil {
		return ResolvedSecret{}, err
	}
	if resp == nil {
		return ResolvedSecret{}, fmt.Errorf("resolve secret response missing")
	}

	return ResolvedSecret{Value: resp.GetValue()}, nil
}

func providerFromProto(provider *secretsv1.SecretProvider) (SecretProvider, error) {
	if provider == nil {
		return SecretProvider{}, fmt.Errorf("provider is required")
	}

	meta := provider.GetMeta()
	if meta == nil {
		return SecretProvider{}, fmt.Errorf("provider meta is required")
	}

	createdAt := meta.GetCreatedAt()
	if createdAt == nil {
		return SecretProvider{}, fmt.Errorf("provider createdAt is required")
	}

	providerType, err := fromProtoProviderType(provider.GetType())
	if err != nil {
		return SecretProvider{}, err
	}

	config, err := providerConfigFromProto(provider.GetConfig())
	if err != nil {
		return SecretProvider{}, err
	}

	if providerType == SecretProviderTypeVault && config.Vault == nil {
		return SecretProvider{}, fmt.Errorf("vault provider missing config")
	}
	if providerType != SecretProviderTypeVault && config.Vault != nil {
		return SecretProvider{}, fmt.Errorf("unsupported provider config")
	}

	var updatedAt *time.Time
	if meta.GetUpdatedAt() != nil {
		parsed := meta.GetUpdatedAt().AsTime().UTC()
		updatedAt = &parsed
	}

	return SecretProvider{
		ID:          meta.GetId(),
		Title:       provider.GetTitle(),
		Description: provider.GetDescription(),
		Type:        providerType,
		Config:      config,
		CreatedAt:   createdAt.AsTime().UTC(),
		UpdatedAt:   updatedAt,
	}, nil
}

func secretFromProto(secret *secretsv1.Secret) (Secret, error) {
	if secret == nil {
		return Secret{}, fmt.Errorf("secret is required")
	}

	meta := secret.GetMeta()
	if meta == nil {
		return Secret{}, fmt.Errorf("secret meta is required")
	}

	createdAt := meta.GetCreatedAt()
	if createdAt == nil {
		return Secret{}, fmt.Errorf("secret createdAt is required")
	}

	var updatedAt *time.Time
	if meta.GetUpdatedAt() != nil {
		parsed := meta.GetUpdatedAt().AsTime().UTC()
		updatedAt = &parsed
	}

	return Secret{
		ID:               meta.GetId(),
		Title:            secret.GetTitle(),
		Description:      secret.GetDescription(),
		SecretProviderID: secret.GetSecretProviderId(),
		RemoteName:       secret.GetRemoteName(),
		CreatedAt:        createdAt.AsTime().UTC(),
		UpdatedAt:        updatedAt,
	}, nil
}

func providerConfigFromProto(config *secretsv1.SecretProviderConfig) (SecretProviderConfig, error) {
	if config == nil {
		return SecretProviderConfig{}, fmt.Errorf("provider config is required")
	}

	vault := config.GetVault()
	if vault == nil {
		return SecretProviderConfig{}, fmt.Errorf("provider config is missing vault")
	}

	return SecretProviderConfig{
		Vault: &VaultConfig{
			Address: vault.GetAddress(),
			Token:   vault.GetToken(),
		},
	}, nil
}

func toProtoProviderConfig(config SecretProviderConfig) (*secretsv1.SecretProviderConfig, error) {
	if config.Vault == nil {
		return nil, fmt.Errorf("provider config must include vault")
	}

	return &secretsv1.SecretProviderConfig{
		Provider: &secretsv1.SecretProviderConfig_Vault{
			Vault: &secretsv1.VaultConfig{
				Address: config.Vault.Address,
				Token:   config.Vault.Token,
			},
		},
	}, nil
}

func toProtoProviderType(providerType SecretProviderType) (secretsv1.SecretProviderType, error) {
	switch providerType {
	case SecretProviderTypeVault:
		return secretsv1.SecretProviderType_SECRET_PROVIDER_TYPE_VAULT, nil
	default:
		return secretsv1.SecretProviderType_SECRET_PROVIDER_TYPE_UNSPECIFIED, fmt.Errorf("unsupported provider type %q", providerType)
	}
}

func fromProtoProviderType(providerType secretsv1.SecretProviderType) (SecretProviderType, error) {
	switch providerType {
	case secretsv1.SecretProviderType_SECRET_PROVIDER_TYPE_VAULT:
		return SecretProviderTypeVault, nil
	default:
		return "", fmt.Errorf("unsupported provider type %q", providerType.String())
	}
}
