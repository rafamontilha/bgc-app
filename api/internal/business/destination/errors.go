package destination

import "errors"

// Erros de validação
var (
	ErrInvalidNCM        = errors.New("NCM deve ter exatamente 8 dígitos numéricos")
	ErrInvalidVolume     = errors.New("volume deve ser um número positivo")
	ErrInvalidMaxResults = errors.New("max_results deve estar entre 1 e 50")
	ErrInvalidCountry    = errors.New("código de país inválido")
)

// Erros de negócio
var (
	ErrNCMNotFound        = errors.New("NCM não encontrado")
	ErrCountryNotFound    = errors.New("país não encontrado")
	ErrNoDataAvailable    = errors.New("dados não disponíveis para o período solicitado")
	ErrInsufficientData   = errors.New("dados insuficientes para gerar recomendações")
)

// Erros de infraestrutura
var (
	ErrDatabaseConnection = errors.New("erro de conexão com banco de dados")
	ErrCacheUnavailable   = errors.New("cache temporariamente indisponível")
	ErrExternalAPIFailed  = errors.New("falha ao consultar API externa")
)
