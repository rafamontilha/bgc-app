package auth

import (
	"fmt"
	"os"
	"path/filepath"
)

// SimpleCertificateManager implementação simples de Certificate Manager
// Para MVP, usa filesystem. Em produção, integraria com vault/k8s secrets
type SimpleCertificateManager struct {
	certsDir string
}

// NewSimpleCertificateManager cria um novo certificate manager
func NewSimpleCertificateManager(certsDir string) *SimpleCertificateManager {
	return &SimpleCertificateManager{
		certsDir: certsDir,
	}
}

// GetCertificate obtém caminho do certificado e chave
func (cm *SimpleCertificateManager) GetCertificate(ref string) (certPath, keyPath string, err error) {
	// Para MVP, certificados estão em certs/{ref}.pem e certs/{ref}.key
	certPath = filepath.Join(cm.certsDir, ref+".pem")
	keyPath = filepath.Join(cm.certsDir, ref+".key")

	// Verifica se arquivos existem
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("certificate not found: %s", ref)
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("certificate key not found: %s", ref)
	}

	return certPath, keyPath, nil
}

// SimpleSecretStore implementação simples de Secret Store
// Para MVP, usa env vars. Em produção, integraria com vault/k8s secrets
type SimpleSecretStore struct{}

// NewSimpleSecretStore cria um novo secret store
func NewSimpleSecretStore() *SimpleSecretStore {
	return &SimpleSecretStore{}
}

// GetSecret obtém secret por referência
func (ss *SimpleSecretStore) GetSecret(ref string) (string, error) {
	// Para MVP, secrets estão em variáveis de ambiente
	// Formato: SECRET_{REF_UPPER}
	// Exemplo: secret-name -> SECRET_SECRET_NAME

	envKey := "SECRET_" + normalizeRef(ref)
	value := os.Getenv(envKey)

	if value == "" {
		return "", fmt.Errorf("secret not found: %s (env: %s)", ref, envKey)
	}

	return value, nil
}

// normalizeRef normaliza referência para env var
func normalizeRef(ref string) string {
	// Remove hífens e converte para uppercase
	result := ""
	for _, c := range ref {
		if c == '-' {
			result += "_"
		} else if c >= 'a' && c <= 'z' {
			result += string(c - 32)
		} else if c >= 'A' && c <= 'Z' {
			result += string(c)
		} else if c >= '0' && c <= '9' {
			result += string(c)
		}
	}
	return result
}
