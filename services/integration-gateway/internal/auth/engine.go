package auth

import (
	"fmt"
	"net/http"

	"github.com/bgc/integration-gateway/internal/types"
)

// Authenticator interface para autenticadores
type Authenticator interface {
	// Authenticate adiciona autenticação a uma requisição HTTP
	Authenticate(req *http.Request) error

	// Type retorna o tipo de autenticação
	Type() string
}

// Engine gerencia autenticadores
type Engine struct {
	certManager CertificateManager
	secretStore SecretStore
}

// CertificateManager interface para gerenciar certificados
type CertificateManager interface {
	GetCertificate(ref string) (certPath, keyPath string, err error)
}

// SecretStore interface para acessar secrets
type SecretStore interface {
	GetSecret(ref string) (string, error)
}

// NewEngine cria um novo auth engine
func NewEngine(certManager CertificateManager, secretStore SecretStore) *Engine {
	return &Engine{
		certManager: certManager,
		secretStore: secretStore,
	}
}

// GetAuthenticator cria o authenticator apropriado baseado na config
func (e *Engine) GetAuthenticator(config *types.AuthConfig) (Authenticator, error) {
	switch config.Type {
	case "none":
		return &NoneAuthenticator{}, nil

	case "api_key":
		if config.APIKey == nil {
			return nil, fmt.Errorf("api_key config is required for api_key auth")
		}

		apiKey, err := e.secretStore.GetSecret(config.APIKey.KeyRef)
		if err != nil {
			return nil, fmt.Errorf("failed to get API key: %w", err)
		}

		return NewAPIKeyAuthenticator(config.APIKey.HeaderName, apiKey), nil

	case "oauth2":
		if config.OAuth2 == nil {
			return nil, fmt.Errorf("oauth2 config is required for oauth2 auth")
		}

		clientSecret, err := e.secretStore.GetSecret(config.OAuth2.ClientSecretRef)
		if err != nil {
			return nil, fmt.Errorf("failed to get OAuth2 client secret: %w", err)
		}

		return NewOAuth2Authenticator(
			config.OAuth2.TokenURL,
			config.OAuth2.ClientID,
			clientSecret,
			config.OAuth2.Scopes,
		), nil

	case "mtls":
		if config.CertificateRef == "" {
			return nil, fmt.Errorf("certificate_ref is required for mtls auth")
		}

		certPath, keyPath, err := e.certManager.GetCertificate(config.CertificateRef)
		if err != nil {
			return nil, fmt.Errorf("failed to get certificate: %w", err)
		}

		return NewMTLSAuthenticator(certPath, keyPath)

	case "basic":
		// TODO: Implementar Basic Auth se necessário
		return nil, fmt.Errorf("basic auth not implemented yet")

	case "jwt":
		// TODO: Implementar JWT se necessário
		return nil, fmt.Errorf("jwt auth not implemented yet")

	default:
		return nil, fmt.Errorf("unknown auth type: %s", config.Type)
	}
}
