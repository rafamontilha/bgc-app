package framework

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bgc/integration-gateway/internal/auth"
	"github.com/bgc/integration-gateway/internal/observability"
	"github.com/bgc/integration-gateway/internal/registry"
	"github.com/bgc/integration-gateway/internal/transform"
	"github.com/bgc/integration-gateway/internal/types"
)

// Executor orquestra a execução de requests usando conectores
type Executor struct {
	registry    *registry.Registry
	authEngine  *auth.Engine
	transformer *transform.Engine
	httpClient  *HTTPClient
}

// NewExecutor cria um novo executor
func NewExecutor(
	reg *registry.Registry,
	authEngine *auth.Engine,
	transformer *transform.Engine,
) *Executor {
	return &Executor{
		registry:    reg,
		authEngine:  authEngine,
		transformer: transformer,
		httpClient:  nil, // Será criado por request (pode variar por config)
	}
}

// Execute executa uma chamada ao connector
func (e *Executor) Execute(ctx *types.ExecutionContext) (*types.ExecutionResult, error) {
	startTime := time.Now()

	// Log início
	observability.WithFields(
		"connector", ctx.ConnectorID,
		"endpoint", ctx.EndpointName,
		"environment", ctx.Environment,
	).Info("Executing connector request")

	// 1. Carrega configuração do connector
	connectorConfig, err := e.registry.Get(ctx.ConnectorID)
	if err != nil {
		observability.RecordError(ctx.ConnectorID, ctx.EndpointName, "connector_not_found")
		return nil, fmt.Errorf("connector not found: %w", err)
	}

	// 2. Obtém configuração do endpoint
	endpointConfig, exists := connectorConfig.Integration.Endpoints[ctx.EndpointName]
	if !exists {
		return nil, fmt.Errorf("endpoint not found: %s", ctx.EndpointName)
	}

	// 3. Obtém URL base do ambiente
	environment, exists := connectorConfig.Environments[ctx.Environment]
	if !exists {
		return nil, fmt.Errorf("environment not found: %s", ctx.Environment)
	}

	// 4. Cria HTTP Client com resiliência
	httpClient := NewHTTPClient(&connectorConfig.Integration.Resilience)

	// 5. Configura autenticação
	authenticator, err := e.authEngine.GetAuthenticator(&connectorConfig.Integration.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticator: %w", err)
	}

	// 6. Constrói URL
	url := e.buildURL(environment.BaseURL, endpointConfig.Path, ctx.Params)

	// 7. Constrói request
	req, err := e.buildRequest(context.Background(), &endpointConfig, url, ctx.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// 8. Aplica autenticação
	if err := authenticator.Authenticate(req); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// 9. Configura TLS para mTLS se necessário
	if mtlsAuth, ok := authenticator.(*auth.MTLSAuthenticator); ok {
		httpClient = e.createMTLSClient(mtlsAuth.GetTLSConfig(), &connectorConfig.Integration.Resilience)
	}

	// 10. Parse timeout
	timeout, _ := parseDuration(endpointConfig.Timeout)
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	// 11. Executa request com resiliência
	resp, err := httpClient.Do(req, timeout)
	if err != nil {
		duration := time.Since(startTime).Seconds()
		observability.RecordRequest(ctx.ConnectorID, ctx.EndpointName, "error", duration)
		observability.RecordError(ctx.ConnectorID, ctx.EndpointName, "http_request_failed")
		observability.WithFields(
			"connector", ctx.ConnectorID,
			"endpoint", ctx.EndpointName,
			"error", err.Error(),
			"duration", duration,
		).Error("HTTP request failed")

		return &types.ExecutionResult{
			Error:    err,
			Duration: time.Since(startTime),
		}, err
	}
	defer resp.Body.Close()

	// 12. Lê response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		duration := time.Since(startTime).Seconds()
		observability.RecordRequest(ctx.ConnectorID, ctx.EndpointName, "error", duration)
		observability.RecordError(ctx.ConnectorID, ctx.EndpointName, "response_read_failed")

		return &types.ExecutionResult{
			StatusCode: resp.StatusCode,
			Error:      fmt.Errorf("failed to read response: %w", err),
			Duration:   time.Since(startTime),
		}, err
	}

	// 13. Verifica status code
	if !e.isSuccessStatus(resp.StatusCode, endpointConfig.Response.SuccessStatus) {
		duration := time.Since(startTime).Seconds()
		observability.RecordRequest(ctx.ConnectorID, ctx.EndpointName, "http_error", duration)
		observability.RecordError(ctx.ConnectorID, ctx.EndpointName, fmt.Sprintf("http_%d", resp.StatusCode))
		observability.WithFields(
			"connector", ctx.ConnectorID,
			"endpoint", ctx.EndpointName,
			"status_code", resp.StatusCode,
			"duration", duration,
		).Warn("Request returned non-success status code")

		return &types.ExecutionResult{
			StatusCode: resp.StatusCode,
			Error:      fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body)),
			Duration:   time.Since(startTime),
		}, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	// 14. Transforma response (JSONPath + plugins)
	data, err := e.transformer.Transform(body, &endpointConfig.Response)
	if err != nil {
		duration := time.Since(startTime).Seconds()
		observability.RecordRequest(ctx.ConnectorID, ctx.EndpointName, "transform_error", duration)
		observability.RecordError(ctx.ConnectorID, ctx.EndpointName, "transform_failed")
		observability.WithFields(
			"connector", ctx.ConnectorID,
			"endpoint", ctx.EndpointName,
			"error", err.Error(),
		).Error("Failed to transform response")

		return &types.ExecutionResult{
			StatusCode: resp.StatusCode,
			Error:      fmt.Errorf("failed to transform response: %w", err),
			Duration:   time.Since(startTime),
		}, err
	}

	// 15. Sucesso! Registra métricas
	duration := time.Since(startTime).Seconds()
	observability.RecordRequest(ctx.ConnectorID, ctx.EndpointName, "success", duration)
	observability.WithFields(
		"connector", ctx.ConnectorID,
		"endpoint", ctx.EndpointName,
		"status_code", resp.StatusCode,
		"duration", duration,
	).Info("Request completed successfully")

	// 16. Retorna resultado
	return &types.ExecutionResult{
		Data:       data,
		StatusCode: resp.StatusCode,
		Duration:   time.Since(startTime),
		Error:      nil,
	}, nil
}

// buildURL constrói URL substituindo placeholders
func (e *Executor) buildURL(baseURL, path string, params map[string]interface{}) string {
	url := baseURL + path

	// Substitui path params {param}
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		url = strings.ReplaceAll(url, placeholder, fmt.Sprintf("%v", value))
	}

	return url
}

// buildRequest constrói HTTP request
func (e *Executor) buildRequest(
	ctx context.Context,
	config *types.EndpointConfig,
	url string,
	params map[string]interface{},
) (*http.Request, error) {
	// Body
	var body io.Reader
	if config.Body != nil && config.Body.Template != "" {
		bodyStr := e.applyTemplate(config.Body.Template, params)
		body = bytes.NewBufferString(bodyStr)
	}

	// Cria request
	req, err := http.NewRequestWithContext(ctx, config.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Query params
	if len(config.QueryParams) > 0 {
		q := req.URL.Query()
		for _, param := range config.QueryParams {
			if value, exists := params[param.Name]; exists {
				q.Add(param.Name, fmt.Sprintf("%v", value))
			} else if param.Default != nil {
				q.Add(param.Name, fmt.Sprintf("%v", param.Default))
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

// applyTemplate substitui {{param}} no template
func (e *Executor) applyTemplate(template string, params map[string]interface{}) string {
	result := template
	for key, value := range params {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// isSuccessStatus verifica se status code é sucesso
func (e *Executor) isSuccessStatus(statusCode int, successStatuses []int) bool {
	for _, status := range successStatuses {
		if statusCode == status {
			return true
		}
	}
	return false
}

// createMTLSClient cria HTTP client com configuração mTLS
func (e *Executor) createMTLSClient(tlsConfig *tls.Config, resilience *types.ResilienceConfig) *HTTPClient {
	baseClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    tlsConfig,
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: false,
		},
	}

	// Cria HTTPClient com resiliência usando o client customizado
	client := NewHTTPClient(resilience)
	client.client = baseClient

	return client
}
