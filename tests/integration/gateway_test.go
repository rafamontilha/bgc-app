package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bgc/integration-gateway/cmd/gateway"
	"github.com/bgc/integration-gateway/internal/auth"
	"github.com/bgc/integration-gateway/internal/framework"
	"github.com/bgc/integration-gateway/internal/registry"
	"github.com/bgc/integration-gateway/internal/transform"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestViaCEPIntegration testa integração completa com mock do ViaCEP
func TestViaCEPIntegration(t *testing.T) {
	// Skip se não estiver em ambiente de CI/CD
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test (set RUN_INTEGRATION_TESTS=true to run)")
	}

	// Mock ViaCEP API
	viacepMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ws/01310100/json/", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"cep": "01310100",
			"logradouro": "Avenida Paulista",
			"complemento": "",
			"bairro": "Bela Vista",
			"localidade": "São Paulo",
			"uf": "SP",
			"ibge": "3550308"
		}`))
	}))
	defer viacepMock.Close()

	// Cria config temporário do ViaCEP apontando para o mock
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "viacep.yaml")

	viacepConfig := `
id: viacep
name: ViaCEP - Consulta de CEP
version: 1.0.0

integration:
  type: rest_api
  auth:
    type: none

  endpoints:
    consulta_cep:
      method: GET
      path: /ws/{cep}/json/
      timeout: 10s

      path_params:
        - name: cep
          type: string
          required: true

      response:
        success_status: [200]
        mapping:
          cep: $.cep
          logradouro: $.logradouro
          bairro: $.bairro
          localidade: $.localidade
          uf: $.uf
        transforms:
          - field: cep
            operation: format_cep

environments:
  development:
    base_url: ` + viacepMock.URL + `
`

	err := os.WriteFile(configPath, []byte(viacepConfig), 0644)
	require.NoError(t, err)

	// Inicializa componentes
	reg := registry.NewRegistry(tmpDir)
	err = reg.LoadAll()
	require.NoError(t, err)

	certManager := auth.NewSimpleCertificateManager(tmpDir)
	secretStore := auth.NewSimpleSecretStore()
	authEngine := auth.NewEngine(certManager, secretStore)

	transformEngine := transform.NewEngine()
	transformEngine.RegisterPlugin("format_cep", &transform.FormatCEPPlugin{})

	executor := framework.NewExecutor(reg, authEngine, transformEngine)

	// Executa request
	ctx := &framework.ExecutionContext{
		ConnectorID:  "viacep",
		EndpointName: "consulta_cep",
		Environment:  "development",
		Params: map[string]interface{}{
			"cep": "01310100",
		},
	}

	result, err := executor.Execute(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotNil(t, result.Data)

	// Verifica transformação do CEP
	assert.Equal(t, "01310-100", result.Data["cep"])
	assert.Equal(t, "Avenida Paulista", result.Data["logradouro"])
	assert.Equal(t, "Bela Vista", result.Data["bairro"])
	assert.Equal(t, "São Paulo", result.Data["localidade"])
	assert.Equal(t, "SP", result.Data["uf"])
}

// TestAPIEndpointsIntegration testa endpoints da API REST
func TestAPIEndpointsIntegration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}

	// Cria config temporário simples
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-api.yaml")

	testConfig := `
id: test-api
name: Test API
version: 1.0.0

integration:
  type: rest_api
  auth:
    type: none
  endpoints:
    get_data:
      method: GET
      path: /data
      response:
        success_status: [200]

environments:
  development:
    base_url: http://test.local
`

	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	require.NoError(t, err)

	// Setup
	reg := registry.NewRegistry(tmpDir)
	err = reg.LoadAll()
	require.NoError(t, err)

	// Cria router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":     "healthy",
			"connectors": reg.Count(),
		})
	})

	// List connectors endpoint
	router.GET("/v1/connectors", func(c *gin.Context) {
		connectors := reg.List()
		result := make([]gin.H, 0, len(connectors))
		for _, conn := range connectors {
			result = append(result, gin.H{
				"id":      conn.ID,
				"name":    conn.Name,
				"version": conn.Version,
			})
		}
		c.JSON(200, result)
	})

	// Test: Health endpoint
	t.Run("Health Endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "healthy", response["status"])
		assert.Equal(t, float64(1), response["connectors"])
	})

	// Test: List connectors endpoint
	t.Run("List Connectors Endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/connectors", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 1)
		assert.Equal(t, "test-api", response[0]["id"])
		assert.Equal(t, "Test API", response[0]["name"])
	})
}

// TestErrorHandling testa tratamento de erros
func TestErrorHandling(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}

	tmpDir := t.TempDir()

	// Inicializa com diretório vazio (sem conectores)
	reg := registry.NewRegistry(tmpDir)
	err := reg.LoadAll()
	require.NoError(t, err)

	certManager := auth.NewSimpleCertificateManager(tmpDir)
	secretStore := auth.NewSimpleSecretStore()
	authEngine := auth.NewEngine(certManager, secretStore)

	transformEngine := transform.NewEngine()
	executor := framework.NewExecutor(reg, authEngine, transformEngine)

	// Test: Connector não encontrado
	ctx := &framework.ExecutionContext{
		ConnectorID:  "nonexistent",
		EndpointName: "test",
		Environment:  "development",
		Params:       map[string]interface{}{},
	}

	_, err = executor.Execute(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connector not found")
}
