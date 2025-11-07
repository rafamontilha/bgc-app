package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bgc/integration-gateway/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader_LoadConnector(t *testing.T) {
	// Cria diretório temporário
	tmpDir := t.TempDir()

	// Cria arquivo de config válido
	validConfig := `
id: test-connector
name: Test Connector
version: 1.0.0
provider: Test Provider

integration:
  type: rest_api
  protocol: https

  auth:
    type: none

  endpoints:
    test_endpoint:
      method: GET
      path: /test
      timeout: 30s
      response:
        success_status: [200]
        mapping:
          id: $.id

environments:
  production:
    base_url: https://api.test.com
`

	configPath := filepath.Join(tmpDir, "test-connector.yaml")
	err := os.WriteFile(configPath, []byte(validConfig), 0644)
	require.NoError(t, err)

	// Testa loader
	loader := NewLoader(tmpDir)
	config, err := loader.LoadConnector("test-connector")

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-connector", config.ID)
	assert.Equal(t, "Test Connector", config.Name)
	assert.Equal(t, "1.0.0", config.Version)
	assert.Equal(t, "rest_api", config.Integration.Type)
	assert.Equal(t, "none", config.Integration.Auth.Type)
	assert.Len(t, config.Integration.Endpoints, 1)
}

func TestLoader_LoadConnector_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewLoader(tmpDir)

	_, err := loader.LoadConnector("nonexistent")
	assert.Error(t, err)
}

func TestLoader_LoadConnector_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	invalidConfig := `
id: invalid
this is not valid yaml: [
`

	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(configPath, []byte(invalidConfig), 0644)
	require.NoError(t, err)

	loader := NewLoader(tmpDir)
	_, err = loader.LoadConnector("invalid")

	assert.Error(t, err)
}

func TestLoader_ValidateConfig(t *testing.T) {
	loader := NewLoader(".")

	tests := []struct {
		name        string
		id          string
		configName  string
		integration map[string]interface{}
		wantErr     bool
		errContains string
	}{
		{
			name:        "Missing ID",
			id:          "",
			configName:  "Test",
			wantErr:     true,
			errContains: "ID is required",
		},
		{
			name:        "Invalid ID format",
			id:          "Invalid_ID",
			configName:  "Test",
			wantErr:     true,
			errContains: "invalid ID format",
		},
		{
			name:       "Valid ID",
			id:         "valid-connector-123",
			configName: "Valid Connector",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.ConnectorConfig{
				ID:      tt.id,
				Name:    tt.configName,
				Version: "1.0.0",
				Integration: types.IntegrationConfig{
					Type: "rest_api",
					Auth: types.AuthConfig{
						Type: "none",
					},
					Endpoints: map[string]types.EndpointConfig{
						"test": {
							Method: "GET",
							Path:   "/test",
							Response: types.ResponseConfig{
								SuccessStatus: []int{200},
							},
						},
					},
				},
			}

			err := loader.validateConfig(config)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoader_ValidateEndpoint(t *testing.T) {
	loader := NewLoader(".")

	tests := []struct {
		name        string
		endpoint    types.EndpointConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid endpoint",
			endpoint: types.EndpointConfig{
				Method: "GET",
				Path:   "/api/test",
				Response: types.ResponseConfig{
					SuccessStatus: []int{200},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid HTTP method",
			endpoint: types.EndpointConfig{
				Method: "INVALID",
				Path:   "/api/test",
			},
			wantErr:     true,
			errContains: "invalid HTTP method",
		},
		{
			name: "Missing path",
			endpoint: types.EndpointConfig{
				Method: "GET",
				Path:   "",
			},
			wantErr:     true,
			errContains: "path is required",
		},
		{
			name: "Path without leading slash",
			endpoint: types.EndpointConfig{
				Method: "GET",
				Path:   "api/test",
			},
			wantErr:     true,
			errContains: "path must start with /",
		},
		{
			name: "Missing success status",
			endpoint: types.EndpointConfig{
				Method: "GET",
				Path:   "/api/test",
				Response: types.ResponseConfig{
					SuccessStatus: []int{},
				},
			},
			wantErr:     true,
			errContains: "at least one success status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loader.validateEndpoint("test_endpoint", &tt.endpoint)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidID(t *testing.T) {
	tests := []struct {
		id      string
		isValid bool
	}{
		{"valid-connector", true},
		{"valid-connector-123", true},
		{"connector123", true},
		{"123-connector", true},
		{"Invalid_ID", false},
		{"Invalid ID", false},
		{"Invalid-ID!", false},
		{"", false},
		{"UPPERCASE", false},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			result := isValidID(tt.id)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestLoader_LoadAllConnectors(t *testing.T) {
	tmpDir := t.TempDir()

	// Cria 2 configs válidos
	configs := []string{
		`
id: connector1
name: Connector 1
version: 1.0.0

integration:
  type: rest_api
  auth:
    type: none
  endpoints:
    test:
      method: GET
      path: /test
      response:
        success_status: [200]

environments:
  production:
    base_url: https://api1.test.com
`,
		`
id: connector2
name: Connector 2
version: 2.0.0

integration:
  type: rest_api
  auth:
    type: api_key
  endpoints:
    fetch:
      method: POST
      path: /fetch
      response:
        success_status: [200, 201]

environments:
  production:
    base_url: https://api2.test.com
`,
	}

	for i, config := range configs {
		path := filepath.Join(tmpDir, filepath.Base(tmpDir)+"-connector"+string(rune('1'+i))+".yaml")
		err := os.WriteFile(path, []byte(config), 0644)
		require.NoError(t, err)
	}

	// Testa LoadAll
	loader := NewLoader(tmpDir)
	connectors, err := loader.LoadAllConnectors()

	assert.NoError(t, err)
	assert.Len(t, connectors, 2)
}
