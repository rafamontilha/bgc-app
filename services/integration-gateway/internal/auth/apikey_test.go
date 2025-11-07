package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyAuthenticator_Authenticate(t *testing.T) {
	auth := NewAPIKeyAuthenticator("X-API-Key", "test-key-123")

	req := httptest.NewRequest("GET", "http://example.com", nil)
	err := auth.Authenticate(req)

	assert.NoError(t, err)
	assert.Equal(t, "test-key-123", req.Header.Get("X-API-Key"))
}

func TestAPIKeyAuthenticator_CustomHeader(t *testing.T) {
	auth := NewAPIKeyAuthenticator("X-Custom-Auth", "custom-key")

	req := httptest.NewRequest("GET", "http://example.com", nil)
	err := auth.Authenticate(req)

	assert.NoError(t, err)
	assert.Equal(t, "custom-key", req.Header.Get("X-Custom-Auth"))
}

func TestAPIKeyAuthenticator_DefaultHeader(t *testing.T) {
	auth := NewAPIKeyAuthenticator("", "key-value")

	req := httptest.NewRequest("GET", "http://example.com", nil)
	err := auth.Authenticate(req)

	assert.NoError(t, err)
	// Deve usar header padrão
	assert.Equal(t, "key-value", req.Header.Get("X-API-Key"))
}

func TestAPIKeyAuthenticator_Type(t *testing.T) {
	auth := NewAPIKeyAuthenticator("X-API-Key", "test")
	assert.Equal(t, "api_key", auth.Type())
}
