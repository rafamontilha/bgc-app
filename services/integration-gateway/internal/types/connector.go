package types

import (
	"time"
)

// ConnectorConfig representa a configuração completa de um connector (mapeado do YAML)
type ConnectorConfig struct {
	ID            string                 `yaml:"id" json:"id"`
	Name          string                 `yaml:"name" json:"name"`
	Version       string                 `yaml:"version" json:"version"`
	Provider      string                 `yaml:"provider" json:"provider"`
	Integration   IntegrationConfig      `yaml:"integration" json:"integration"`
	Environments  map[string]Environment `yaml:"environments" json:"environments"`
	Compliance    ComplianceConfig       `yaml:"compliance" json:"compliance"`
	Governance    GovernanceConfig       `yaml:"governance" json:"governance"`
	Observability ObservabilityConfig    `yaml:"observability" json:"observability"`
}

// IntegrationConfig configuração de integração
type IntegrationConfig struct {
	Type       string                    `yaml:"type" json:"type"`     // rest_api, soap, graphql, grpc
	Protocol   string                    `yaml:"protocol" json:"protocol"` // http, https
	Auth       AuthConfig                `yaml:"auth" json:"auth"`
	Endpoints  map[string]EndpointConfig `yaml:"endpoints" json:"endpoints"`
	Resilience ResilienceConfig          `yaml:"resilience" json:"resilience"`
	Cache      CacheConfig               `yaml:"cache" json:"cache"`
}

// AuthConfig configuração de autenticação
type AuthConfig struct {
	Type           string         `yaml:"type" json:"type"` // mtls, oauth2, api_key, basic, jwt, none
	CertificateRef string         `yaml:"certificate_ref,omitempty" json:"certificate_ref,omitempty"`
	OAuth2         *OAuth2Config  `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
	APIKey         *APIKeyConfig  `yaml:"api_key,omitempty" json:"api_key,omitempty"`
}

// OAuth2Config configuração OAuth2
type OAuth2Config struct {
	TokenURL        string   `yaml:"token_url" json:"token_url"`
	ClientID        string   `yaml:"client_id" json:"client_id"`
	ClientSecretRef string   `yaml:"client_secret_ref" json:"client_secret_ref"`
	Scopes          []string `yaml:"scopes" json:"scopes"`
}

// APIKeyConfig configuração de API Key
type APIKeyConfig struct {
	HeaderName string `yaml:"header_name" json:"header_name"`
	KeyRef     string `yaml:"key_ref" json:"key_ref"`
}

// EndpointConfig configuração de um endpoint
type EndpointConfig struct {
	Method      string                 `yaml:"method" json:"method"` // GET, POST, PUT, PATCH, DELETE
	Path        string                 `yaml:"path" json:"path"`
	Timeout     string                 `yaml:"timeout" json:"timeout"`
	Headers     map[string]string      `yaml:"headers" json:"headers"`
	QueryParams []ParameterConfig      `yaml:"query_params,omitempty" json:"query_params,omitempty"`
	PathParams  []ParameterConfig      `yaml:"path_params,omitempty" json:"path_params,omitempty"`
	Body        *BodyConfig            `yaml:"body,omitempty" json:"body,omitempty"`
	Response    ResponseConfig         `yaml:"response" json:"response"`
	Plugins     map[string]string      `yaml:"plugins,omitempty" json:"plugins,omitempty"`
}

// ParameterConfig configuração de parâmetro
type ParameterConfig struct {
	Name      string      `yaml:"name" json:"name"`
	Type      string      `yaml:"type" json:"type"` // string, integer, number, boolean
	Required  bool        `yaml:"required" json:"required"`
	Format    string      `yaml:"format,omitempty" json:"format,omitempty"`
	Pattern   string      `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MinLength int         `yaml:"min_length,omitempty" json:"min_length,omitempty"`
	MaxLength int         `yaml:"max_length,omitempty" json:"max_length,omitempty"`
	Default   interface{} `yaml:"default,omitempty" json:"default,omitempty"`
}

// BodyConfig configuração de body
type BodyConfig struct {
	ContentType string `yaml:"content_type" json:"content_type"`
	Template    string `yaml:"template" json:"template"`
}

// ResponseConfig configuração de resposta
type ResponseConfig struct {
	SuccessStatus []int                  `yaml:"success_status" json:"success_status"`
	ErrorStatus   []int                  `yaml:"error_status" json:"error_status"`
	Mapping       map[string]string      `yaml:"mapping" json:"mapping"` // field -> JSONPath
	Transforms    []TransformConfig      `yaml:"transforms,omitempty" json:"transforms,omitempty"`
}

// TransformConfig configuração de transformação
type TransformConfig struct {
	Field     string                 `yaml:"field" json:"field"`
	Operation string                 `yaml:"operation" json:"operation"`
	Values    map[string]string      `yaml:"values,omitempty" json:"values,omitempty"`
	Params    map[string]interface{} `yaml:"params,omitempty" json:"params,omitempty"`
}

// ResilienceConfig configuração de resiliência
type ResilienceConfig struct {
	Retry          *RetryConfig          `yaml:"retry,omitempty" json:"retry,omitempty"`
	CircuitBreaker *CircuitBreakerConfig `yaml:"circuit_breaker,omitempty" json:"circuit_breaker,omitempty"`
	RateLimit      *RateLimitConfig      `yaml:"rate_limit,omitempty" json:"rate_limit,omitempty"`
}

// RetryConfig configuração de retry
type RetryConfig struct {
	MaxAttempts     int    `yaml:"max_attempts" json:"max_attempts"`
	Backoff         string `yaml:"backoff" json:"backoff"` // constant, linear, exponential
	InitialInterval string `yaml:"initial_interval" json:"initial_interval"`
	MaxInterval     string `yaml:"max_interval" json:"max_interval"`
}

// CircuitBreakerConfig configuração de circuit breaker
type CircuitBreakerConfig struct {
	FailureThreshold int    `yaml:"failure_threshold" json:"failure_threshold"`
	SuccessThreshold int    `yaml:"success_threshold" json:"success_threshold"`
	Timeout          string `yaml:"timeout" json:"timeout"`
}

// RateLimitConfig configuração de rate limit
type RateLimitConfig struct {
	RequestsPerMinute int `yaml:"requests_per_minute" json:"requests_per_minute"`
	Burst             int `yaml:"burst" json:"burst"`
}

// CacheConfig configuração de cache
type CacheConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	TTL        string `yaml:"ttl" json:"ttl"`
	KeyPattern string `yaml:"key_pattern" json:"key_pattern"`
}

// Environment configuração de ambiente
type Environment struct {
	BaseURL     string `yaml:"base_url" json:"base_url"`
	HealthCheck string `yaml:"health_check" json:"health_check"`
}

// ComplianceConfig configuração de compliance
type ComplianceConfig struct {
	Tags               []string `yaml:"tags" json:"tags"`
	DataClassification string   `yaml:"data_classification" json:"data_classification"`
	RetentionDays      int      `yaml:"retention_days" json:"retention_days"`
	EncryptionRequired bool     `yaml:"encryption_required" json:"encryption_required"`
}

// GovernanceConfig configuração de governança
type GovernanceConfig struct {
	OwnerTeam       string `yaml:"owner_team" json:"owner_team"`
	ApprovedBy      string `yaml:"approved_by" json:"approved_by"`
	LastAudited     string `yaml:"last_audited" json:"last_audited"`
	ReviewFrequency string `yaml:"review_frequency" json:"review_frequency"`
}

// ObservabilityConfig configuração de observabilidade
type ObservabilityConfig struct {
	MetricsEnabled bool          `yaml:"metrics_enabled" json:"metrics_enabled"`
	TracingEnabled bool          `yaml:"tracing_enabled" json:"tracing_enabled"`
	LogLevel       string        `yaml:"log_level" json:"log_level"`
	Alerts         []AlertConfig `yaml:"alerts" json:"alerts"`
}

// AlertConfig configuração de alerta
type AlertConfig struct {
	Type      string   `yaml:"type" json:"type"` // certificate_expiry, error_rate, latency, availability
	Threshold string   `yaml:"threshold" json:"threshold"`
	Window    string   `yaml:"window,omitempty" json:"window,omitempty"`
	Channels  []string `yaml:"channels" json:"channels"`
}

// ExecutionContext contexto de execução de uma requisição
type ExecutionContext struct {
	ConnectorID  string
	EndpointName string
	Environment  string
	Params       map[string]interface{}
	StartTime    time.Time
}

// ExecutionResult resultado da execução
type ExecutionResult struct {
	Data       map[string]interface{}
	StatusCode int
	Duration   time.Duration
	Error      error
	CacheHit   bool
	RetryCount int
}
