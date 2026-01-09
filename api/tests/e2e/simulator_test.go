// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://api.bgc.local"

// TestSimulatorHappyPath testa request mínimo
func TestSimulatorHappyPath(t *testing.T) {
	payload := map[string]interface{}{
		"ncm": "17011400",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	destinations, ok := result["destinations"].([]interface{})
	require.True(t, ok, "destinations should be an array")
	assert.GreaterOrEqual(t, len(destinations), 6, "Should return at least 6 destinations")

	// Verificar primeiro destino tem todos os campos
	first := destinations[0].(map[string]interface{})
	assert.NotEmpty(t, first["country_code"])
	assert.NotEmpty(t, first["country_name"])
	assert.NotZero(t, first["score"])
	assert.NotZero(t, first["rank"])
	assert.NotEmpty(t, first["demand_level"])
	assert.NotEmpty(t, first["recommendation_reason"])
}

// TestSimulatorWithCountryFilter testa filtro de países
func TestSimulatorWithCountryFilter(t *testing.T) {
	payload := map[string]interface{}{
		"ncm":       "17011400",
		"countries": []string{"US", "CN", "DE"},
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	destinations, _ := result["destinations"].([]interface{})
	assert.LessOrEqual(t, len(destinations), 3, "Should return max 3 destinations")

	// Verificar que apenas US, CN, DE estão presentes
	for _, dest := range destinations {
		d := dest.(map[string]interface{})
		code := d["country_code"].(string)
		assert.Contains(t, []string{"US", "CN", "DE"}, code)
	}
}

// TestSimulatorWithVolume testa request com volume
func TestSimulatorWithVolume(t *testing.T) {
	payload := map[string]interface{}{
		"ncm":         "26011200",
		"volume_kg":   5000,
		"max_results": 5,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	destinations, _ := result["destinations"].([]interface{})
	assert.LessOrEqual(t, len(destinations), 5, "Should return at most 5 destinations")

	// Verificar que logistics_cost_usd foi calculado com volume
	if len(destinations) > 0 {
		first := destinations[0].(map[string]interface{})
		assert.NotZero(t, first["logistics_cost_usd"])
	}
}

// TestSimulatorInvalidNCM testa NCM inválido
func TestSimulatorInvalidNCM(t *testing.T) {
	payload := map[string]interface{}{
		"ncm": "12345", // Apenas 5 dígitos
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 400, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "validation_error", result["error"])
	assert.Contains(t, result["message"], "NCM")
}

// TestSimulatorNCMNotFound testa NCM não encontrado
func TestSimulatorNCMNotFound(t *testing.T) {
	payload := map[string]interface{}{
		"ncm": "99999999", // NCM inexistente
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Pode retornar 404 ou 200 com array vazio dependendo da implementação
	assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 200)

	if resp.StatusCode == 404 {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "ncm_not_found", result["error"])
	} else {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		destinations := result["destinations"].([]interface{})
		assert.Equal(t, 0, len(destinations))
	}
}

// TestSimulatorRateLimiting testa rate limiting
func TestSimulatorRateLimiting(t *testing.T) {
	payload := map[string]interface{}{
		"ncm": "17011400",
	}
	body, _ := json.Marshal(payload)

	client := &http.Client{}

	// Fazer 5 requests (limite free)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("POST", baseURL+"/v1/simulator/destinations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verificar headers de rate limit
		assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Limit"))
		assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Remaining"))

		resp.Body.Close()
		time.Sleep(100 * time.Millisecond)
	}

	// 6ª request deve ser bloqueada
	req, _ := http.NewRequest("POST", baseURL+"/v1/simulator/destinations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 429, resp.StatusCode)

	// Verificar headers
	assert.Equal(t, "5", resp.Header.Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", resp.Header.Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Reset"))
}

// TestSimulatorPerformance testa performance da API
func TestSimulatorPerformance(t *testing.T) {
	payload := map[string]interface{}{
		"ncm": "17011400",
	}
	body, _ := json.Marshal(payload)

	start := time.Now()
	resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer(body))
	duration := time.Since(start)

	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	// Performance target: < 200ms (current: 2-4ms)
	assert.Less(t, duration.Milliseconds(), int64(200), "Response time should be < 200ms")

	t.Logf("Response time: %v ms", duration.Milliseconds())
}
