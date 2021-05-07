package oidc

import (
	"github.com/ory/fosite"
	"gopkg.in/square/go-jose.v2"

	"github.com/authelia/authelia/internal/authentication"
	"github.com/authelia/authelia/internal/authorization"
	"github.com/authelia/authelia/internal/configuration/schema"
	"github.com/authelia/authelia/internal/session"
)

// NewClient creates a new InternalClient.
func NewClient(config schema.OpenIDConnectClientConfiguration) *InternalClient {
	return &InternalClient{
		ID:            config.ID,
		Description:   config.Description,
		Policy:        authorization.PolicyToLevel(config.Policy),
		Secret:        []byte(config.Secret),
		RedirectURIs:  config.RedirectURIs,
		GrantTypes:    config.GrantTypes,
		ResponseTypes: config.ResponseTypes,
		Scopes:        config.Scopes,

		ResponseModes: []fosite.ResponseModeType{
			fosite.ResponseModeDefault,
			fosite.ResponseModeFormPost,
			fosite.ResponseModeQuery,
			fosite.ResponseModeFragment,
		},
		TokenEndpointAuthMethod: "client_secret_post",
	}
}

// InternalClient represents the client internally.
type InternalClient struct {
	ID            string              `json:"id"`
	Description   string              `json:"-"`
	Secret        []byte              `json:"client_secret,omitempty"`
	RedirectURIs  []string            `json:"redirect_uris"`
	GrantTypes    []string            `json:"grant_types"`
	ResponseTypes []string            `json:"response_types"`
	Scopes        []string            `json:"scopes"`
	Audience      []string            `json:"audience"`
	Public        bool                `json:"public"`
	Policy        authorization.Level `json:"-"`

	// These are the OpenIDConnect Client props.
	RequestURIs                       []string                  `json:"request_uris"`
	ResponseModes                     []fosite.ResponseModeType `json:"response_modes"`
	JSONWebKeys                       *jose.JSONWebKeySet       `json:"jwks,omitempty"`
	JSONWebKeysURI                    string                    `json:"jwks_uri,omitempty"`
	RequestObjectSigningAlgorithm     string                    `json:"request_object_signing_alg,omitempty"`
	TokenEndpointAuthSigningAlgorithm string                    `json:"token_endpoint_auth_signing_alg,omitempty"`
	TokenEndpointAuthMethod           string                    `json:"token_endpoint_auth_method,omitempty"`
}

// IsAuthenticationLevelSufficient returns if the provided authentication.Level is sufficient for the client of the AutheliaClient.
func (c InternalClient) IsAuthenticationLevelSufficient(level authentication.Level) bool {
	return authorization.IsAuthLevelSufficient(level, c.Policy)
}

// GetID returns the ID.
func (c InternalClient) GetID() string {
	return c.ID
}

func (c InternalClient) GetConsentRequestBody(session *session.OIDCWorkflowSession) (body ConsentGetResponseBody) {
	body.ClientID = c.ID
	body.ClientDescription = c.Description
	if session != nil {
		body.Scopes = scopeNamesToScopes(session.RequestedScopes)
		body.Audience = audienceNamesToAudience(session.RequestedAudience)
	}

	return body
}

// GetHashedSecret returns the Secret.
func (c InternalClient) GetHashedSecret() []byte {
	return c.Secret
}

// GetRedirectURIs returns the RedirectURIs.
func (c InternalClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

// GetGrantTypes returns the GrantTypes.
func (c InternalClient) GetGrantTypes() fosite.Arguments {
	if len(c.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}

	return c.GrantTypes
}

// GetResponseTypes returns the ResponseTypes.
func (c InternalClient) GetResponseTypes() fosite.Arguments {
	if len(c.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}

	return c.ResponseTypes
}

// GetScopes returns the Scopes.
func (c InternalClient) GetScopes() fosite.Arguments {
	return c.Scopes
}

// IsPublic returns the value of the Public property.
func (c InternalClient) IsPublic() bool {
	return c.Public
}

// GetAudience returns the Audience.
func (c InternalClient) GetAudience() fosite.Arguments {
	return c.Audience
}

// GetResponseModes returns the valid response modes for this client.
//
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetResponseModes() []fosite.ResponseModeType {
	return c.ResponseModes
}

// GetRequestURIs returns the valid request URIs for this client.
//
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetRequestURIs() []string {
	return c.RequestURIs
}

// GetJSONWebKeys returns the jose.JSONWebKeySet for this client.
//
// This should only be set for RP's that don't understand the JSONWebKeysURI, and should not be set along side this.
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetJSONWebKeys() *jose.JSONWebKeySet {
	return c.JSONWebKeys
}

// GetJSONWebKeysURI returns the URI containing the JWKS for this client.
//
// This should generally always be set unless the client does not understand the JWKS URI in which case
// GetJSONWebKeys should be set instead.
//
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetJSONWebKeysURI() string {
	return c.JSONWebKeysURI
}

// GetRequestObjectSigningAlgorithm returns the object sig alg.
//
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetRequestObjectSigningAlgorithm() string {
	return c.RequestObjectSigningAlgorithm
}

// GetTokenEndpointAuthMethod returns the token endpoint auth method.
//
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetTokenEndpointAuthMethod() string {
	if c.TokenEndpointAuthMethod == "" {
		return "client_secret_basic"
	}

	return c.TokenEndpointAuthMethod
}

// GetTokenEndpointAuthSigningAlgorithm returns the token auth sig alg.
//
// Implements part of fosite.OpenIDConnectClient.
func (c InternalClient) GetTokenEndpointAuthSigningAlgorithm() string {
	if c.TokenEndpointAuthSigningAlgorithm == "" {
		return "RS256"
	}

	return c.TokenEndpointAuthSigningAlgorithm
}
