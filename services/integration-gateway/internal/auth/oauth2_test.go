package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOAuth2Authenticator_Authenticate(t *testing.T) {
	// Mock OAuth2 server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/token", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Retorna token válido
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "test-token-123",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer server.Close()

	auth := NewOAuth2Authenticator(
		server.URL+"/token",
		"test-client-id",
		"test-client-secret",
		[]string{"read", "write"},
	)

	req := httptest.NewRequest("GET", "http://example.com", nil)

	err := auth.Authenticate(req)
	assert.NoError(t, err)

	// Verifica header Authorization
	authHeader := req.Header.Get("Authorization")
	assert.Equal(t, "Bearer test-token-123", authHeader)
}

func TestOAuth2Authenticator_TokenCaching(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "cached-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer server.Close()

	auth := NewOAuth2Authenticator(
		server.URL+"/token",
		"test-client",
		"test-secret",
		nil,
	)

	// Primeira chamada - deve fazer request
	req1 := httptest.NewRequest("GET", "http://example.com", nil)
	err := auth.Authenticate(req1)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Segunda chamada - deve usar cache
	req2 := httptest.NewRequest("GET", "http://example.com", nil)
	err = auth.Authenticate(req2)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount, "Should use cached token")

	// Verifica que ambos têm o mesmo token
	assert.Equal(t, req1.Header.Get("Authorization"), req2.Header.Get("Authorization"))
}

func TestOAuth2Authenticator_TokenExpiration(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Retorna token que expira em 1 segundo
		w.Write([]byte(`{
			"access_token": "expiring-token",
			"token_type": "Bearer",
			"expires_in": 1
		}`))
	}))
	defer server.Close()

	auth := NewOAuth2Authenticator(
		server.URL+"/token",
		"test-client",
		"test-secret",
		nil,
	)

	// Primeira chamada
	req1 := httptest.NewRequest("GET", "http://example.com", nil)
	err := auth.Authenticate(req1)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Aguarda expiração (token expira em 1s, mas renova 5min antes = -299s)
	// Como o tempo de renovação é negativo, deve renovar imediatamente
	time.Sleep(10 * time.Millisecond)

	// Segunda chamada - deve renovar token
	req2 := httptest.NewRequest("GET", "http://example.com", nil)
	err = auth.Authenticate(req2)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount, "Should renew expired token")
}

func TestOAuth2Authenticator_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid_client"}`))
	}))
	defer server.Close()

	auth := NewOAuth2Authenticator(
		server.URL+"/token",
		"invalid-client",
		"invalid-secret",
		nil,
	)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	err := auth.Authenticate(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func TestOAuth2Authenticator_Type(t *testing.T) {
	auth := NewOAuth2Authenticator("http://test.com/token", "id", "secret", nil)
	assert.Equal(t, "oauth2", auth.Type())
}
