package registry

import (
	"fmt"
	"sync"

	"github.com/bgc/integration-gateway/internal/types"
)

// Registry gerencia os conectores carregados
type Registry struct {
	mu         sync.RWMutex
	connectors map[string]*types.ConnectorConfig
	loader     *Loader
}

// NewRegistry cria um novo registry
func NewRegistry(configDir string) *Registry {
	return &Registry{
		connectors: make(map[string]*types.ConnectorConfig),
		loader:     NewLoader(configDir),
	}
}

// LoadAll carrega todos os conectores do diretório
func (r *Registry) LoadAll() error {
	configs, err := r.loader.LoadAllConnectors()
	if err != nil {
		return fmt.Errorf("failed to load connectors: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, config := range configs {
		r.connectors[config.ID] = config
	}

	return nil
}

// Get obtém um connector pelo ID
func (r *Registry) Get(id string) (*types.ConnectorConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	config, exists := r.connectors[id]
	if !exists {
		return nil, fmt.Errorf("connector not found: %s", id)
	}

	return config, nil
}

// List retorna todos os conectores registrados
func (r *Registry) List() []*types.ConnectorConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*types.ConnectorConfig, 0, len(r.connectors))
	for _, config := range r.connectors {
		result = append(result, config)
	}

	return result
}

// Reload recarrega um connector específico
func (r *Registry) Reload(id string) error {
	config, err := r.loader.LoadConnector(id)
	if err != nil {
		return fmt.Errorf("failed to reload connector %s: %w", id, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.connectors[id] = config

	return nil
}

// Count retorna o número de conectores registrados
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.connectors)
}

// GetEndpoint obtém configuração de um endpoint específico
func (r *Registry) GetEndpoint(connectorID, endpointName string) (*types.EndpointConfig, error) {
	config, err := r.Get(connectorID)
	if err != nil {
		return nil, err
	}

	endpoint, exists := config.Integration.Endpoints[endpointName]
	if !exists {
		return nil, fmt.Errorf("endpoint not found: %s.%s", connectorID, endpointName)
	}

	return &endpoint, nil
}

// GetEnvironment obtém configuração de ambiente
func (r *Registry) GetEnvironment(connectorID, env string) (*types.Environment, error) {
	config, err := r.Get(connectorID)
	if err != nil {
		return nil, err
	}

	environment, exists := config.Environments[env]
	if !exists {
		return nil, fmt.Errorf("environment not found: %s for connector %s", env, connectorID)
	}

	return &environment, nil
}
