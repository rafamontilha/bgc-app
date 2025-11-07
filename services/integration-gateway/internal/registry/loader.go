package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bgc/integration-gateway/internal/types"
	"gopkg.in/yaml.v3"
)

// Loader carrega configurações de conectores de arquivos YAML
type Loader struct {
	configDir string
}

// NewLoader cria um novo loader
func NewLoader(configDir string) *Loader {
	return &Loader{
		configDir: configDir,
	}
}

// LoadConnector carrega um connector específico pelo ID
func (l *Loader) LoadConnector(id string) (*types.ConnectorConfig, error) {
	filename := filepath.Join(l.configDir, id+".yaml")

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read connector config %s: %w", id, err)
	}

	var config types.ConnectorConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse connector config %s: %w", id, err)
	}

	// Validação básica
	if err := l.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid connector config %s: %w", id, err)
	}

	return &config, nil
}

// LoadAllConnectors carrega todos os conectores do diretório
func (l *Loader) LoadAllConnectors() ([]*types.ConnectorConfig, error) {
	files, err := filepath.Glob(filepath.Join(l.configDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list connector configs: %w", err)
	}

	var configs []*types.ConnectorConfig
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", file, err)
		}

		var config types.ConnectorConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}

		if err := l.validateConfig(&config); err != nil {
			return nil, fmt.Errorf("invalid config in %s: %w", file, err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

// validateConfig valida a configuração do connector
func (l *Loader) validateConfig(config *types.ConnectorConfig) error {
	if config.ID == "" {
		return fmt.Errorf("connector ID is required")
	}

	if config.Name == "" {
		return fmt.Errorf("connector name is required")
	}

	if config.Version == "" {
		return fmt.Errorf("connector version is required")
	}

	// Valida formato do ID (lowercase, hyphens only)
	if !isValidID(config.ID) {
		return fmt.Errorf("invalid ID format: must be lowercase with hyphens only")
	}

	// Valida tipo de integração
	validTypes := map[string]bool{
		"rest_api": true,
		"soap":     true,
		"graphql":  true,
		"grpc":     true,
	}
	if !validTypes[config.Integration.Type] {
		return fmt.Errorf("invalid integration type: %s", config.Integration.Type)
	}

	// Valida tipo de auth
	validAuthTypes := map[string]bool{
		"mtls":   true,
		"oauth2": true,
		"api_key": true,
		"basic":  true,
		"jwt":    true,
		"none":   true,
	}
	if !validAuthTypes[config.Integration.Auth.Type] {
		return fmt.Errorf("invalid auth type: %s", config.Integration.Auth.Type)
	}

	// Valida que há pelo menos um endpoint
	if len(config.Integration.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint is required")
	}

	// Valida endpoints
	for name, endpoint := range config.Integration.Endpoints {
		if err := l.validateEndpoint(name, &endpoint); err != nil {
			return err
		}
	}

	return nil
}

// validateEndpoint valida configuração de endpoint
func (l *Loader) validateEndpoint(name string, endpoint *types.EndpointConfig) error {
	// Valida HTTP method
	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"PATCH":  true,
		"DELETE": true,
	}
	if !validMethods[endpoint.Method] {
		return fmt.Errorf("invalid HTTP method for endpoint %s: %s", name, endpoint.Method)
	}

	// Valida path
	if endpoint.Path == "" {
		return fmt.Errorf("path is required for endpoint %s", name)
	}

	if !strings.HasPrefix(endpoint.Path, "/") {
		return fmt.Errorf("path must start with / for endpoint %s", name)
	}

	// Valida response
	if len(endpoint.Response.SuccessStatus) == 0 {
		return fmt.Errorf("at least one success status is required for endpoint %s", name)
	}

	return nil
}

// isValidID verifica se o ID está no formato correto
func isValidID(id string) bool {
	if id == "" {
		return false
	}

	for _, c := range id {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}

	return true
}
