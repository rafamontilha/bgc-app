package auth

import "net/http"

// APIKeyAuthenticator autenticação via API Key em header
type APIKeyAuthenticator struct {
	headerName string
	apiKey     string
}

// NewAPIKeyAuthenticator cria um novo authenticator de API Key
func NewAPIKeyAuthenticator(headerName, apiKey string) *APIKeyAuthenticator {
	if headerName == "" {
		headerName = "X-API-Key"
	}

	return &APIKeyAuthenticator{
		headerName: headerName,
		apiKey:     apiKey,
	}
}

func (a *APIKeyAuthenticator) Authenticate(req *http.Request) error {
	req.Header.Set(a.headerName, a.apiKey)
	return nil
}

func (a *APIKeyAuthenticator) Type() string {
	return "api_key"
}
