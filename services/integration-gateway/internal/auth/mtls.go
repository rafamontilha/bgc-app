package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// MTLSAuthenticator autenticação via mTLS (mutual TLS)
type MTLSAuthenticator struct {
	certPath string
	keyPath  string
	tlsConfig *tls.Config
}

// NewMTLSAuthenticator cria um novo authenticator mTLS
func NewMTLSAuthenticator(certPath, keyPath string) (*MTLSAuthenticator, error) {
	// Carrega certificado e chave
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	// Carrega CA certificates (ICP-Brasil chain)
	caCertPool, err := loadCACertPool()
	if err != nil {
		// Log warning mas não falha (usa system CA pool)
		caCertPool = nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
		// Para produção, validar sempre. Para dev/sandbox, pode ser necessário:
		// InsecureSkipVerify: os.Getenv("MTLS_SKIP_VERIFY") == "true",
	}

	return &MTLSAuthenticator{
		certPath:  certPath,
		keyPath:   keyPath,
		tlsConfig: tlsConfig,
	}, nil
}

func (a *MTLSAuthenticator) Authenticate(req *http.Request) error {
	// mTLS é configurado no transport do http.Client, não na request
	// O executor precisa usar o transport com TLS config

	// Verifica se o client do request tem transport personalizado
	if req != nil && req.Header != nil {
		// Adiciona headers adicionais se necessário para ICP-Brasil
		req.Header.Set("X-Client-Cert", "present")
	}

	return nil
}

func (a *MTLSAuthenticator) Type() string {
	return "mtls"
}

// GetTLSConfig retorna a configuração TLS para usar no http.Client
func (a *MTLSAuthenticator) GetTLSConfig() *tls.Config {
	return a.tlsConfig
}

// loadCACertPool carrega CA certificates (ICP-Brasil)
func loadCACertPool() (*x509.CertPool, error) {
	// Tenta carregar de arquivo configurado
	caCertPath := os.Getenv("ICP_CA_CERT_PATH")
	if caCertPath == "" {
		// Usa system CA pool
		return nil, nil
	}

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA cert")
	}

	return caCertPool, nil
}
