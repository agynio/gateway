package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/agynio/gateway/internal/llmclient"
	"github.com/agynio/gateway/internal/llmgen"
)

func providerFromClient(provider llmclient.LLMProvider) (llmgen.LLMProvider, error) {
	id, err := uuidFromString(provider.ID)
	if err != nil {
		return llmgen.LLMProvider{}, fmt.Errorf("parse provider id: %w", err)
	}

	authMethod, err := authMethodFromClient(provider.AuthMethod)
	if err != nil {
		return llmgen.LLMProvider{}, err
	}

	createdAt := provider.CreatedAt.UTC()
	var updatedAt *time.Time
	if provider.UpdatedAt != nil {
		value := provider.UpdatedAt.UTC()
		updatedAt = &value
	}

	return llmgen.LLMProvider{
		Id:         id,
		Endpoint:   provider.Endpoint,
		AuthMethod: authMethod,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func modelFromClient(model llmclient.Model) (llmgen.Model, error) {
	id, err := uuidFromString(model.ID)
	if err != nil {
		return llmgen.Model{}, fmt.Errorf("parse model id: %w", err)
	}

	providerID, err := uuidFromString(model.LLMProviderID)
	if err != nil {
		return llmgen.Model{}, fmt.Errorf("parse model provider id: %w", err)
	}

	createdAt := model.CreatedAt.UTC()
	var updatedAt *time.Time
	if model.UpdatedAt != nil {
		value := model.UpdatedAt.UTC()
		updatedAt = &value
	}

	return llmgen.Model{
		Id:            id,
		LlmProviderId: providerID,
		Name:          model.Name,
		RemoteName:    model.RemoteName,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func providerCreateToParams(body llmgen.LLMProviderCreateRequest) (llmclient.CreateProviderParams, error) {
	authMethod, err := authMethodToClient(body.AuthMethod)
	if err != nil {
		return llmclient.CreateProviderParams{}, err
	}

	return llmclient.CreateProviderParams{
		Endpoint:   body.Endpoint,
		AuthMethod: authMethod,
		Token:      body.Token,
	}, nil
}

func providerUpdateToParams(body llmgen.LLMProviderUpdateRequest) (llmclient.UpdateProviderParams, error) {
	params := llmclient.UpdateProviderParams{
		Endpoint: body.Endpoint,
		Token:    body.Token,
	}

	if body.AuthMethod != nil {
		authMethod, err := authMethodToClient(*body.AuthMethod)
		if err != nil {
			return llmclient.UpdateProviderParams{}, err
		}
		params.AuthMethod = &authMethod
	}

	return params, nil
}

func modelCreateToParams(body llmgen.ModelCreateRequest) (llmclient.CreateModelParams, error) {
	return llmclient.CreateModelParams{
		Name:          body.Name,
		LLMProviderID: body.LlmProviderId.String(),
		RemoteName:    body.RemoteName,
	}, nil
}

func modelUpdateToParams(body llmgen.ModelUpdateRequest) (llmclient.UpdateModelParams, error) {
	params := llmclient.UpdateModelParams{
		Name:       body.Name,
		RemoteName: body.RemoteName,
	}

	if body.LlmProviderId != nil {
		providerID := body.LlmProviderId.String()
		params.LLMProviderID = &providerID
	}

	return params, nil
}

func authMethodFromClient(method llmclient.AuthMethod) (llmgen.AuthMethod, error) {
	switch method {
	case llmclient.AuthMethodBearer:
		return llmgen.Bearer, nil
	default:
		return "", fmt.Errorf("unsupported auth method %q", method)
	}
}

func authMethodToClient(method llmgen.AuthMethod) (llmclient.AuthMethod, error) {
	switch method {
	case llmgen.Bearer:
		return llmclient.AuthMethodBearer, nil
	default:
		return "", fmt.Errorf("unsupported auth method %q", method)
	}
}

func uuidFromString(value string) (openapi_types.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("parse uuid: %w", err)
	}
	return openapi_types.UUID(parsed), nil
}
