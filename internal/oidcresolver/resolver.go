package oidcresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	usersv1 "github.com/agynio/gateway/gen/agynio/api/users/v1"
	"github.com/agynio/gateway/internal/identity"
	"github.com/agynio/gateway/internal/oidcauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProfileSource selects where provisioning-time profile claims are read from.
type ProfileSource string

const (
	// ProfileSourceUserInfo fetches claims from the IdP's UserInfo endpoint.
	ProfileSourceUserInfo ProfileSource = "userinfo"
	// ProfileSourceToken reads claims from the validated access token itself.
	// Required for IdPs that issue audience-restricted access tokens: the OIDC
	// UserInfo endpoint rejects any token carrying an `aud` claim, so a Gateway
	// that needs a verifiable JWT cannot also call UserInfo with it.
	ProfileSourceToken ProfileSource = "token"
)

// ClaimNames maps profile fields onto the claims that carry them. Applies to
// both sources, so a provider using non-standard names is configured the same
// way regardless of where claims are read from.
type ClaimNames struct {
	Name              string
	Email             string
	Picture           string
	PreferredUsername string
}

// DefaultClaimNames are the standard OIDC claim names.
func DefaultClaimNames() ClaimNames {
	return ClaimNames{
		Name:              "name",
		Email:             "email",
		Picture:           "picture",
		PreferredUsername: "preferred_username",
	}
}

func (c ClaimNames) withDefaults() ClaimNames {
	defaults := DefaultClaimNames()
	if strings.TrimSpace(c.Name) == "" {
		c.Name = defaults.Name
	}
	if strings.TrimSpace(c.Email) == "" {
		c.Email = defaults.Email
	}
	if strings.TrimSpace(c.Picture) == "" {
		c.Picture = defaults.Picture
	}
	if strings.TrimSpace(c.PreferredUsername) == "" {
		c.PreferredUsername = defaults.PreferredUsername
	}
	return c
}

// Option customizes how a Resolver obtains profile claims.
type Option func(*Resolver)

// WithProfileSource selects the claim source. Defaults to ProfileSourceUserInfo.
func WithProfileSource(source ProfileSource) Option {
	return func(r *Resolver) {
		r.profileSource = source
	}
}

// WithClaimNames overrides the claim names. Empty fields keep their defaults.
func WithClaimNames(names ClaimNames) Option {
	return func(r *Resolver) {
		r.claimNames = names.withDefaults()
	}
}

type Resolver struct {
	verifier         *oidcauth.Verifier
	usersClient      usersv1.UsersServiceClient
	userinfoEndpoint string
	httpClient       *http.Client
	profileSource    ProfileSource
	claimNames       ClaimNames
}

func NewResolver(verifier *oidcauth.Verifier, usersClient usersv1.UsersServiceClient, httpClient *http.Client, opts ...Option) (*Resolver, error) {
	if verifier == nil {
		return nil, fmt.Errorf("verifier is required")
	}
	if usersClient == nil {
		return nil, fmt.Errorf("users client is required")
	}
	if httpClient == nil {
		return nil, fmt.Errorf("http client is required")
	}

	resolver := &Resolver{
		verifier:         verifier,
		usersClient:      usersClient,
		userinfoEndpoint: strings.TrimSpace(verifier.UserinfoEndpoint()),
		httpClient:       httpClient,
		profileSource:    ProfileSourceUserInfo,
		claimNames:       DefaultClaimNames(),
	}
	for _, opt := range opts {
		opt(resolver)
	}

	switch resolver.profileSource {
	case ProfileSourceToken:
	case ProfileSourceUserInfo:
		// Only the UserInfo source needs the endpoint, so a provider that omits
		// it from discovery is still usable with the token source.
		if resolver.userinfoEndpoint == "" {
			return nil, fmt.Errorf("userinfo endpoint is required")
		}
	default:
		return nil, fmt.Errorf("unsupported profile source %q", resolver.profileSource)
	}

	return resolver, nil
}

func (r *Resolver) ResolveFromToken(ctx context.Context, accessToken string) (identity.ResolvedIdentity, error) {
	claims, err := r.verifier.Verify(ctx, accessToken)
	if err != nil {
		return identity.ResolvedIdentity{}, status.Errorf(codes.Unauthenticated, "invalid bearer token: %v", err)
	}

	getResponse, err := r.usersClient.GetUserByOIDCSubject(ctx, &usersv1.GetUserByOIDCSubjectRequest{
		OidcSubject: claims.Subject,
	})
	if err == nil {
		return identityFromUser(getResponse.GetUser())
	}
	if status.Code(err) != codes.NotFound {
		return identity.ResolvedIdentity{}, err
	}

	profileClaims, err := r.profileClaims(ctx, accessToken, claims)
	if err != nil {
		return identity.ResolvedIdentity{}, err
	}

	preferredUsername := stringClaim(profileClaims, r.claimNames.PreferredUsername)
	var preferredUsernamePtr *string
	if preferredUsername != "" {
		preferredUsernamePtr = &preferredUsername
	}

	createResponse, err := r.usersClient.ResolveOrCreateUser(ctx, &usersv1.ResolveOrCreateUserRequest{
		OidcSubject:       claims.Subject,
		Name:              stringClaim(profileClaims, r.claimNames.Name),
		Email:             stringClaim(profileClaims, r.claimNames.Email),
		PhotoUrl:          stringClaim(profileClaims, r.claimNames.Picture),
		PreferredUsername: preferredUsernamePtr,
	})
	if err != nil {
		return identity.ResolvedIdentity{}, err
	}

	return identityFromUser(createResponse.GetUser())
}

// profileClaims returns the claims used to provision a new user, from whichever
// source is configured.
func (r *Resolver) profileClaims(ctx context.Context, accessToken string, claims oidcauth.Claims) (map[string]any, error) {
	if r.profileSource == ProfileSourceToken {
		// The token's signature, issuer, expiry and subject are already verified.
		return claims.All, nil
	}

	userInfo, err := r.fetchUserInfo(ctx, accessToken, claims.Subject)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user info: %v", err)
	}
	return userInfo, nil
}

// stringClaim reads a claim as a trimmed string. A missing claim, or one that
// isn't a string, yields "" so a partial profile provisions rather than fails.
func stringClaim(claims map[string]any, name string) string {
	if claims == nil || name == "" {
		return ""
	}
	value, ok := claims[name].(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

func (r *Resolver) fetchUserInfo(ctx context.Context, accessToken, expectedSub string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.userinfoEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	const maxUserInfoResponseBytes = 1 << 20
	limitedBody := io.LimitReader(resp.Body, maxUserInfoResponseBytes)

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(limitedBody)
		if readErr != nil {
			return nil, fmt.Errorf("userinfo request failed with status %s", resp.Status)
		}
		return nil, fmt.Errorf("userinfo request failed with status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	// Decoded as a map rather than oidc.UserInfo so configured claim names apply
	// to this source too.
	var userInfo map[string]any
	decoder := json.NewDecoder(limitedBody)
	if err := decoder.Decode(&userInfo); err != nil {
		return nil, err
	}

	subject := stringClaim(userInfo, "sub")
	if subject == "" {
		return nil, fmt.Errorf("userinfo subject is required")
	}
	if subject != expectedSub {
		return nil, fmt.Errorf("userinfo subject does not match expected subject")
	}

	return userInfo, nil
}

func identityFromUser(user *usersv1.User) (identity.ResolvedIdentity, error) {
	if user == nil {
		return identity.ResolvedIdentity{}, fmt.Errorf("user missing from response")
	}
	meta := user.GetMeta()
	if meta == nil {
		return identity.ResolvedIdentity{}, fmt.Errorf("user metadata missing from response")
	}
	identityID := strings.TrimSpace(meta.GetId())
	if identityID == "" {
		return identity.ResolvedIdentity{}, fmt.Errorf("identity id missing")
	}

	return identity.ResolvedIdentity{
		IdentityID:   identityID,
		IdentityType: identity.IdentityTypeUser,
	}, nil
}
