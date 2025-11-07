package framework

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/bgc/integration-gateway/internal/types"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

// HTTPClient cliente HTTP genérico com resiliência
type HTTPClient struct {
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	rateLimiter    *rate.Limiter
	retryConfig    *types.RetryConfig
}

// NewHTTPClient cria um novo cliente HTTP com resiliência
func NewHTTPClient(resilience *types.ResilienceConfig) *HTTPClient {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: false,
		},
	}

	hc := &HTTPClient{
		client: client,
	}

	// Configura Circuit Breaker
	if resilience != nil && resilience.CircuitBreaker != nil {
		cbConfig := resilience.CircuitBreaker
		timeout, _ := parseDuration(cbConfig.Timeout)

		settings := gobreaker.Settings{
			Name:        "http-client",
			MaxRequests: uint32(cbConfig.SuccessThreshold),
			Interval:    timeout,
			Timeout:     timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= 0.6
			},
		}

		hc.circuitBreaker = gobreaker.NewCircuitBreaker(settings)
	}

	// Configura Rate Limiter
	if resilience != nil && resilience.RateLimit != nil {
		rps := float64(resilience.RateLimit.RequestsPerMinute) / 60.0
		burst := resilience.RateLimit.Burst
		if burst == 0 {
			burst = resilience.RateLimit.RequestsPerMinute
		}
		hc.rateLimiter = rate.NewLimiter(rate.Limit(rps), burst)
	}

	// Configura Retry
	if resilience != nil && resilience.Retry != nil {
		hc.retryConfig = resilience.Retry
	}

	return hc
}

// Do executa uma requisição HTTP com resiliência
func (c *HTTPClient) Do(req *http.Request, timeout time.Duration) (*http.Response, error) {
	ctx := req.Context()

	// Apply timeout
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// Rate limiting
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit exceeded: %w", err)
		}
	}

	// Circuit breaker + retry
	if c.circuitBreaker != nil {
		result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
			return c.doWithRetry(req)
		})
		if err != nil {
			return nil, err
		}
		return result.(*http.Response), nil
	}

	return c.doWithRetry(req)
}

// doWithRetry executa requisição com retry
func (c *HTTPClient) doWithRetry(req *http.Request) (*http.Response, error) {
	if c.retryConfig == nil {
		return c.client.Do(req)
	}

	var lastErr error
	maxAttempts := c.retryConfig.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = 1
	}

	initialInterval, _ := parseDuration(c.retryConfig.InitialInterval)
	if initialInterval == 0 {
		initialInterval = 1 * time.Second
	}

	maxInterval, _ := parseDuration(c.retryConfig.MaxInterval)
	if maxInterval == 0 {
		maxInterval = 30 * time.Second
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Clone request for retry (body can only be read once)
		reqClone := req.Clone(req.Context())

		resp, err := c.client.Do(reqClone)

		// Success
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Store error
		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			resp.Body.Close()
		}

		// Don't retry on last attempt
		if attempt == maxAttempts {
			break
		}

		// Calculate backoff
		waitTime := c.calculateBackoff(attempt, initialInterval, maxInterval)

		// Wait before retry
		select {
		case <-time.After(waitTime):
			// Continue to next attempt
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// calculateBackoff calcula tempo de espera baseado na estratégia
func (c *HTTPClient) calculateBackoff(attempt int, initial, max time.Duration) time.Duration {
	var wait time.Duration

	switch c.retryConfig.Backoff {
	case "constant":
		wait = initial

	case "linear":
		wait = time.Duration(attempt) * initial

	case "exponential":
		wait = time.Duration(math.Pow(2, float64(attempt-1))) * initial

	default:
		wait = initial
	}

	if wait > max {
		wait = max
	}

	return wait
}

// parseDuration parse string de duração (ex: "30s", "5m", "1h")
func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	return time.ParseDuration(s)
}

// RequestBuilder helper para construir requisições
type RequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    io.Reader
}

// NewRequestBuilder cria um novo builder
func NewRequestBuilder(method, url string) *RequestBuilder {
	return &RequestBuilder{
		method:  method,
		url:     url,
		headers: make(map[string]string),
	}
}

// WithHeader adiciona header
func (b *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	b.headers[key] = value
	return b
}

// WithHeaders adiciona múltiplos headers
func (b *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	for k, v := range headers {
		b.headers[k] = v
	}
	return b
}

// WithBody adiciona body
func (b *RequestBuilder) WithBody(body io.Reader) *RequestBuilder {
	b.body = body
	return b
}

// Build constrói a requisição
func (b *RequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, b.method, b.url, b.body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
