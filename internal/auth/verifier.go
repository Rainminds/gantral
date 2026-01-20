package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
)

// IdentityType distinguishes betweens humans and machines
type IdentityType string

const (
	IdentityTypeHuman   IdentityType = "human"
	IdentityTypeMachine IdentityType = "machine"
)

// Identity represents an authenticated entity (User or Machine).
type Identity struct {
	Subject  string       `json:"sub"`
	Type     IdentityType `json:"type"`
	Provider string       `json:"provider"`
	Roles    []string     `json:"roles,omitempty"`
	OrgID    string       `json:"org_id,omitempty"`
}

// Claims (internal) standard OIDC claims we parse initially
type internalClaims struct {
	Subject string   `json:"sub"`
	Roles   []string `json:"roles,omitempty"`
	OrgID   string   `json:"org_id,omitempty"`
}

// TokenVerifier validates a raw token string and returns the identity.
type TokenVerifier interface {
	Verify(ctx context.Context, tokenString string) (*Identity, error)
}

// MultiVerifier chains multiple verifiers
type MultiVerifier struct {
	verifiers []TokenVerifier
}

// NewMultiVerifier creates a new verifier that tries each verifier in order.
func NewMultiVerifier(verifiers ...TokenVerifier) *MultiVerifier {
	return &MultiVerifier{verifiers: verifiers}
}

func (mv *MultiVerifier) Verify(ctx context.Context, tokenString string) (*Identity, error) {
	var errs []string
	for _, v := range mv.verifiers {
		identity, err := v.Verify(ctx, tokenString)
		if err == nil {
			return identity, nil
		}
		errs = append(errs, err.Error())
	}
	return nil, fmt.Errorf("authentication failed: %s", strings.Join(errs, "; "))
}

// OIDCVerifier implements TokenVerifier using an OpenID Connect provider.
// It enforces RS256 signatures from the provider.
type OIDCVerifier struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	identityType IdentityType
	providerName string
}

func NewOIDCVerifier(ctx context.Context, issuerURL, clientID string, iType IdentityType) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %w", err)
	}

	conf := oidc.Config{
		ClientID: clientID,
		// Ensure the algorithms are what we expect (default is RS256)
		SupportedSigningAlgs: []string{oidc.RS256},
	}
	verifier := provider.Verifier(&conf)

	return &OIDCVerifier{
		provider:     provider,
		verifier:     verifier,
		identityType: iType,
		providerName: issuerURL,
	}, nil
}

func (v *OIDCVerifier) Verify(ctx context.Context, tokenString string) (*Identity, error) {
	idToken, err := v.verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("oidc verification failed: %w", err)
	}

	var claims internalClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &Identity{
		Subject:  claims.Subject,
		Type:     v.identityType,
		Provider: v.providerName,
		Roles:    claims.Roles,
		OrgID:    claims.OrgID,
	}, nil
}

// DevVerifier implements TokenVerifier for development purposes.
// It accepts HS256 tokens signed with a shared secret.
type DevVerifier struct {
	Secret       []byte
	IdentityType IdentityType
}

func (v *DevVerifier) Verify(ctx context.Context, tokenString string) (*Identity, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt parse failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	mapClaims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	// Manual mapping since MapClaims is weird with struct tags
	sub, _ := (*mapClaims)["sub"].(string)
	orgID, _ := (*mapClaims)["org_id"].(string)

	// Handle roles
	var roles []string
	if r, ok := (*mapClaims)["roles"].([]interface{}); ok {
		for _, role := range r {
			if strRole, ok := role.(string); ok {
				roles = append(roles, strRole)
			}
		}
	}

	if sub == "" {
		return nil, errors.New("token missing subject")
	}

	return &Identity{
		Subject:  sub,
		Type:     v.IdentityType,
		Provider: "dev-mode",
		Roles:    roles,
		OrgID:    orgID,
	}, nil
}
