package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// OAuth2Authenticator autenticação via OAuth2 Client Credentials
type OAuth2Authenticator struct {
	tokenURL     string
	clientID     string
	clientSecret string
	scopes       []string

	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time
	httpClient  *http.Client
}

// OAuth2TokenResponse resposta do token endpoint
type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// NewOAuth2Authenticator cria um novo authenticator OAuth2
func NewOAuth2Authenticator(tokenURL, clientID, clientSecret string, scopes []string) *OAuth2Authenticator {
	return &OAuth2Authenticator{
		tokenURL:     tokenURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       scopes,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (a *OAuth2Authenticator) Authenticate(req *http.Request) error {
	token, err := a.getToken()
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (a *OAuth2Authenticator) Type() string {
	return "oauth2"
}

// getToken obtém token (usa cache se válido)
func (a *OAuth2Authenticator) getToken() (string, error) {
	// Verifica se tem token válido em cache
	a.mu.RLock()
	if a.accessToken != "" && time.Now().Before(a.expiresAt) {
		token := a.accessToken
		a.mu.RUnlock()
		return token, nil
	}
	a.mu.RUnlock()

	// Precisa obter novo token
	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-check (outro goroutine pode ter atualizado)
	if a.accessToken != "" && time.Now().Before(a.expiresAt) {
		return a.accessToken, nil
	}

	// Faz requisição para obter token
	token, expiresIn, err := a.fetchToken()
	if err != nil {
		return "", err
	}

	// Atualiza cache (renova 5min antes de expirar)
	a.accessToken = token
	a.expiresAt = time.Now().Add(time.Duration(expiresIn-300) * time.Second)

	return token, nil
}

// fetchToken faz requisição para obter novo token
func (a *OAuth2Authenticator) fetchToken() (string, int, error) {
	// Prepara body (application/x-www-form-urlencoded)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)
	if len(a.scopes) > 0 {
		data.Set("scope", strings.Join(a.scopes, " "))
	}

	// Cria requisição
	req, err := http.NewRequest("POST", a.tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Executa requisição
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to fetch token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse resposta
	var tokenResp OAuth2TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", 0, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", 0, fmt.Errorf("empty access token in response")
	}

	return tokenResp.AccessToken, tokenResp.ExpiresIn, nil
}
