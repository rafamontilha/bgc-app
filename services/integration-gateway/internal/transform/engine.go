package transform

import (
	"fmt"
	"strings"

	"github.com/bgc/integration-gateway/internal/types"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

// Engine mecanismo de transformação de dados
type Engine struct {
	plugins map[string]TransformPlugin
}

// TransformPlugin interface para plugins de transformação
type TransformPlugin interface {
	Transform(value interface{}, params map[string]interface{}) (interface{}, error)
}

// NewEngine cria um novo engine de transformação
func NewEngine() *Engine {
	return &Engine{
		plugins: make(map[string]TransformPlugin),
	}
}

// RegisterPlugin registra um plugin de transformação
func (e *Engine) RegisterPlugin(name string, plugin TransformPlugin) {
	e.plugins[name] = plugin
}

// Transform aplica transformações nos dados
func (e *Engine) Transform(data interface{}, config *types.ResponseConfig) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Converte data para JSON parseable
	jsonData, err := e.ensureJSONData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	// Aplica mapeamentos JSONPath
	for field, jsonPath := range config.Mapping {
		value, err := e.extractValue(jsonData, jsonPath)
		if err != nil {
			// Log warning mas não falha (campo pode ser opcional)
			continue
		}
		result[field] = value
	}

	// Aplica transformações
	for _, transform := range config.Transforms {
		if value, exists := result[transform.Field]; exists {
			transformed, err := e.applyTransform(value, &transform)
			if err != nil {
				return nil, fmt.Errorf("failed to transform field %s: %w", transform.Field, err)
			}
			result[transform.Field] = transformed
		}
	}

	return result, nil
}

// extractValue extrai valor usando JSONPath
func (e *Engine) extractValue(data interface{}, jsonPath string) (interface{}, error) {
	// Parse JSONPath
	path, err := jp.ParseString(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("invalid JSONPath %s: %w", jsonPath, err)
	}

	// Extract value
	results := path.Get(data)
	if len(results) == 0 {
		return nil, fmt.Errorf("no value found for JSONPath: %s", jsonPath)
	}

	// Se retornou múltiplos valores, retorna array
	if len(results) > 1 {
		return results, nil
	}

	return results[0], nil
}

// applyTransform aplica transformação em um valor
func (e *Engine) applyTransform(value interface{}, config *types.TransformConfig) (interface{}, error) {
	// Transformação map_values (mapeamento de valores)
	if config.Operation == "map_values" {
		return e.mapValues(value, config.Values)
	}

	// Usa plugin registrado
	if plugin, exists := e.plugins[config.Operation]; exists {
		return plugin.Transform(value, config.Params)
	}

	return nil, fmt.Errorf("unknown transform operation: %s", config.Operation)
}

// mapValues mapeia valores baseado em dicionário
func (e *Engine) mapValues(value interface{}, mapping map[string]string) (interface{}, error) {
	strValue := fmt.Sprintf("%v", value)
	if mapped, exists := mapping[strValue]; exists {
		return mapped, nil
	}
	// Se não encontrou mapeamento, retorna valor original
	return value, nil
}

// ensureJSONData garante que data está em formato parseable
func (e *Engine) ensureJSONData(data interface{}) (interface{}, error) {
	// Se já é map ou slice, retorna direto
	switch data.(type) {
	case map[string]interface{}, []interface{}:
		return data, nil
	}

	// Se é string, tenta fazer parse
	if str, ok := data.(string); ok {
		parsed, err := oj.ParseString(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON string: %w", err)
		}
		return parsed, nil
	}

	// Se é []byte, tenta fazer parse
	if bytes, ok := data.([]byte); ok {
		parsed, err := oj.Parse(bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON bytes: %w", err)
		}
		return parsed, nil
	}

	// Retorna como está
	return data, nil
}

// --- Built-in Transform Plugins ---

// FormatCNPJPlugin formata CNPJ (12345678000195 -> 12.345.678/0001-95)
type FormatCNPJPlugin struct{}

func (p *FormatCNPJPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)

	// Remove non-digits
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, str)

	if len(digits) != 14 {
		return value, nil // Retorna original se não for CNPJ válido
	}

	// Format: XX.XXX.XXX/XXXX-XX
	return fmt.Sprintf("%s.%s.%s/%s-%s",
		digits[0:2],
		digits[2:5],
		digits[5:8],
		digits[8:12],
		digits[12:14],
	), nil
}

// FormatCPFPlugin formata CPF (12345678901 -> 123.456.789-01)
type FormatCPFPlugin struct{}

func (p *FormatCPFPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)

	// Remove non-digits
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, str)

	if len(digits) != 11 {
		return value, nil
	}

	// Format: XXX.XXX.XXX-XX
	return fmt.Sprintf("%s.%s.%s-%s",
		digits[0:3],
		digits[3:6],
		digits[6:9],
		digits[9:11],
	), nil
}

// FormatCEPPlugin formata CEP (01310100 -> 01310-100)
type FormatCEPPlugin struct{}

func (p *FormatCEPPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)

	// Remove non-digits
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, str)

	if len(digits) != 8 {
		return value, nil
	}

	// Format: XXXXX-XXX
	return fmt.Sprintf("%s-%s", digits[0:5], digits[5:8]), nil
}

// ToUpperPlugin converte para maiúsculas
type ToUpperPlugin struct{}

func (p *ToUpperPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	return strings.ToUpper(str), nil
}

// ToLowerPlugin converte para minúsculas
type ToLowerPlugin struct{}

func (p *ToLowerPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	return strings.ToLower(str), nil
}

// TrimPlugin remove espaços
type TrimPlugin struct{}

func (p *TrimPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	return strings.TrimSpace(str), nil
}
