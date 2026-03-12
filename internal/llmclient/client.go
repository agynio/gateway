package llmclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	llmv1 "github.com/agynio/gateway/gen/agynio/api/llm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthMethod string

const (
	AuthMethodBearer AuthMethod = "bearer"
)

type Client struct {
	conn   *grpc.ClientConn
	client llmv1.LLMServiceClient
}

type LLMProvider struct {
	ID         string     `json:"id"`
	Endpoint   string     `json:"endpoint"`
	AuthMethod AuthMethod `json:"authMethod"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type Model struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	LLMProviderID string     `json:"llmProviderId"`
	RemoteName    string     `json:"remoteName"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

type CreateProviderParams struct {
	Endpoint   string     `json:"endpoint"`
	AuthMethod AuthMethod `json:"authMethod"`
	Token      string     `json:"token"`
}

type UpdateProviderParams struct {
	Endpoint   *string     `json:"endpoint"`
	AuthMethod *AuthMethod `json:"authMethod"`
	Token      *string     `json:"token"`
}

type CreateModelParams struct {
	Name          string `json:"name"`
	LLMProviderID string `json:"llmProviderId"`
	RemoteName    string `json:"remoteName"`
}

type UpdateModelParams struct {
	Name          *string `json:"name"`
	LLMProviderID *string `json:"llmProviderId"`
	RemoteName    *string `json:"remoteName"`
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
		client: llmv1.NewLLMServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateProvider(ctx context.Context, params CreateProviderParams) (LLMProvider, error) {
	authMethod, err := toProtoAuthMethod(params.AuthMethod)
	if err != nil {
		return LLMProvider{}, err
	}

	resp, err := c.client.CreateLLMProvider(ctx, &llmv1.CreateLLMProviderRequest{
		Endpoint:   params.Endpoint,
		AuthMethod: authMethod,
		Token:      params.Token,
	})
	if err != nil {
		return LLMProvider{}, err
	}

	provider := resp.GetProvider()
	if provider == nil {
		return LLMProvider{}, fmt.Errorf("create provider response missing provider")
	}

	return providerFromProto(provider)
}

func (c *Client) GetProvider(ctx context.Context, id string) (LLMProvider, error) {
	resp, err := c.client.GetLLMProvider(ctx, &llmv1.GetLLMProviderRequest{Id: id})
	if err != nil {
		return LLMProvider{}, err
	}

	provider := resp.GetProvider()
	if provider == nil {
		return LLMProvider{}, fmt.Errorf("get provider response missing provider")
	}

	return providerFromProto(provider)
}

func (c *Client) UpdateProvider(ctx context.Context, id string, params UpdateProviderParams) (LLMProvider, error) {
	req := &llmv1.UpdateLLMProviderRequest{Id: id}
	if params.Endpoint != nil {
		req.Endpoint = params.Endpoint
	}
	if params.AuthMethod != nil {
		method, err := toProtoAuthMethod(*params.AuthMethod)
		if err != nil {
			return LLMProvider{}, err
		}
		req.AuthMethod = &method
	}
	if params.Token != nil {
		req.Token = params.Token
	}

	resp, err := c.client.UpdateLLMProvider(ctx, req)
	if err != nil {
		return LLMProvider{}, err
	}

	provider := resp.GetProvider()
	if provider == nil {
		return LLMProvider{}, fmt.Errorf("update provider response missing provider")
	}

	return providerFromProto(provider)
}

func (c *Client) DeleteProvider(ctx context.Context, id string) error {
	_, err := c.client.DeleteLLMProvider(ctx, &llmv1.DeleteLLMProviderRequest{Id: id})
	return err
}

func (c *Client) ListProviders(ctx context.Context, pageSize int32, pageToken string) ([]LLMProvider, string, error) {
	resp, err := c.client.ListLLMProviders(ctx, &llmv1.ListLLMProvidersRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	})
	if err != nil {
		return nil, "", err
	}

	providers := make([]LLMProvider, 0, len(resp.GetProviders()))
	for _, provider := range resp.GetProviders() {
		converted, err := providerFromProto(provider)
		if err != nil {
			return nil, "", err
		}
		providers = append(providers, converted)
	}
	if providers == nil {
		providers = []LLMProvider{}
	}

	return providers, resp.GetNextPageToken(), nil
}

func (c *Client) CreateModel(ctx context.Context, params CreateModelParams) (Model, error) {
	resp, err := c.client.CreateModel(ctx, &llmv1.CreateModelRequest{
		Name:          params.Name,
		LlmProviderId: params.LLMProviderID,
		RemoteName:    params.RemoteName,
	})
	if err != nil {
		return Model{}, err
	}

	model := resp.GetModel()
	if model == nil {
		return Model{}, fmt.Errorf("create model response missing model")
	}

	return modelFromProto(model)
}

func (c *Client) GetModel(ctx context.Context, id string) (Model, error) {
	resp, err := c.client.GetModel(ctx, &llmv1.GetModelRequest{Id: id})
	if err != nil {
		return Model{}, err
	}

	model := resp.GetModel()
	if model == nil {
		return Model{}, fmt.Errorf("get model response missing model")
	}

	return modelFromProto(model)
}

func (c *Client) UpdateModel(ctx context.Context, id string, params UpdateModelParams) (Model, error) {
	req := &llmv1.UpdateModelRequest{Id: id}
	if params.Name != nil {
		req.Name = params.Name
	}
	if params.LLMProviderID != nil {
		req.LlmProviderId = params.LLMProviderID
	}
	if params.RemoteName != nil {
		req.RemoteName = params.RemoteName
	}

	resp, err := c.client.UpdateModel(ctx, req)
	if err != nil {
		return Model{}, err
	}

	model := resp.GetModel()
	if model == nil {
		return Model{}, fmt.Errorf("update model response missing model")
	}

	return modelFromProto(model)
}

func (c *Client) DeleteModel(ctx context.Context, id string) error {
	_, err := c.client.DeleteModel(ctx, &llmv1.DeleteModelRequest{Id: id})
	return err
}

func (c *Client) ListModels(ctx context.Context, pageSize int32, pageToken, llmProviderID string) ([]Model, string, error) {
	req := &llmv1.ListModelsRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}
	if strings.TrimSpace(llmProviderID) != "" {
		req.LlmProviderId = llmProviderID
	}

	resp, err := c.client.ListModels(ctx, req)
	if err != nil {
		return nil, "", err
	}

	models := make([]Model, 0, len(resp.GetModels()))
	for _, model := range resp.GetModels() {
		converted, err := modelFromProto(model)
		if err != nil {
			return nil, "", err
		}
		models = append(models, converted)
	}
	if models == nil {
		models = []Model{}
	}

	return models, resp.GetNextPageToken(), nil
}

func providerFromProto(provider *llmv1.LLMProvider) (LLMProvider, error) {
	if provider == nil {
		return LLMProvider{}, fmt.Errorf("provider is required")
	}

	meta := provider.GetMeta()
	if meta == nil {
		return LLMProvider{}, fmt.Errorf("provider meta is required")
	}

	createdAt := meta.GetCreatedAt()
	if createdAt == nil {
		return LLMProvider{}, fmt.Errorf("provider createdAt is required")
	}

	authMethod, err := fromProtoAuthMethod(provider.GetAuthMethod())
	if err != nil {
		return LLMProvider{}, err
	}

	var updatedAt *time.Time
	if meta.GetUpdatedAt() != nil {
		parsed := meta.GetUpdatedAt().AsTime().UTC()
		updatedAt = &parsed
	}

	return LLMProvider{
		ID:         meta.GetId(),
		Endpoint:   provider.GetEndpoint(),
		AuthMethod: authMethod,
		CreatedAt:  createdAt.AsTime().UTC(),
		UpdatedAt:  updatedAt,
	}, nil
}

func modelFromProto(model *llmv1.Model) (Model, error) {
	if model == nil {
		return Model{}, fmt.Errorf("model is required")
	}

	meta := model.GetMeta()
	if meta == nil {
		return Model{}, fmt.Errorf("model meta is required")
	}

	createdAt := meta.GetCreatedAt()
	if createdAt == nil {
		return Model{}, fmt.Errorf("model createdAt is required")
	}

	var updatedAt *time.Time
	if meta.GetUpdatedAt() != nil {
		parsed := meta.GetUpdatedAt().AsTime().UTC()
		updatedAt = &parsed
	}

	return Model{
		ID:            meta.GetId(),
		Name:          model.GetName(),
		LLMProviderID: model.GetLlmProviderId(),
		RemoteName:    model.GetRemoteName(),
		CreatedAt:     createdAt.AsTime().UTC(),
		UpdatedAt:     updatedAt,
	}, nil
}

func toProtoAuthMethod(method AuthMethod) (llmv1.AuthMethod, error) {
	switch method {
	case AuthMethodBearer:
		return llmv1.AuthMethod_AUTH_METHOD_BEARER, nil
	default:
		return llmv1.AuthMethod_AUTH_METHOD_UNSPECIFIED, fmt.Errorf("unsupported auth method %q", method)
	}
}

func fromProtoAuthMethod(method llmv1.AuthMethod) (AuthMethod, error) {
	switch method {
	case llmv1.AuthMethod_AUTH_METHOD_BEARER:
		return AuthMethodBearer, nil
	default:
		return "", fmt.Errorf("unsupported auth method %q", method.String())
	}
}
