package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleCertificateManager_GetCertificate(t *testing.T) {
	tmpDir := t.TempDir()

	// Cria certificados de teste
	certPath := filepath.Join(tmpDir, "test-cert.pem")
	keyPath := filepath.Join(tmpDir, "test-cert.key")

	err := os.WriteFile(certPath, []byte("fake-cert-data"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(keyPath, []byte("fake-key-data"), 0644)
	require.NoError(t, err)

	// Testa certificate manager
	cm := NewSimpleCertificateManager(tmpDir)
	cert, key, err := cm.GetCertificate("test-cert")

	assert.NoError(t, err)
	assert.Equal(t, certPath, cert)
	assert.Equal(t, keyPath, key)
}

func TestSimpleCertificateManager_CertificateNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewSimpleCertificateManager(tmpDir)

	_, _, err := cm.GetCertificate("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSimpleCertificateManager_KeyNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Cria apenas certificado, sem chave
	certPath := filepath.Join(tmpDir, "test-cert.pem")
	err := os.WriteFile(certPath, []byte("fake-cert-data"), 0644)
	require.NoError(t, err)

	cm := NewSimpleCertificateManager(tmpDir)
	_, _, err = cm.GetCertificate("test-cert")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
}

func TestSimpleSecretStore_GetSecret(t *testing.T) {
	// Seta variável de ambiente
	os.Setenv("SECRET_TEST_SECRET", "secret-value-123")
	defer os.Unsetenv("SECRET_TEST_SECRET")

	store := NewSimpleSecretStore()
	secret, err := store.GetSecret("test-secret")

	assert.NoError(t, err)
	assert.Equal(t, "secret-value-123", secret)
}

func TestSimpleSecretStore_SecretNotFound(t *testing.T) {
	store := NewSimpleSecretStore()
	_, err := store.GetSecret("nonexistent-secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestNormalizeRef(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test-secret", "TEST_SECRET"},
		{"my-api-key", "MY_API_KEY"},
		{"oauth-client-secret", "OAUTH_CLIENT_SECRET"},
		{"simple", "SIMPLE"},
		{"ALREADY-UPPER", "ALREADY_UPPER"},
		{"mixed-Case-123", "MIXED_CASE_123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeRef(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
