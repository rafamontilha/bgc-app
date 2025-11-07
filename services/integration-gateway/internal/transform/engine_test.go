package transform

import (
	"testing"

	"github.com/bgc/integration-gateway/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestEngine_Transform_JSONPath(t *testing.T) {
	engine := NewEngine()

	// Registra plugins
	engine.RegisterPlugin("format_cnpj", &FormatCNPJPlugin{})

	// JSON de teste
	jsonData := map[string]interface{}{
		"data": map[string]interface{}{
			"cnpj":        "12345678000195",
			"nome":        "Empresa Teste LTDA",
			"situacao":    "01",
			"endereco": map[string]interface{}{
				"logradouro": "Rua Teste",
				"numero":     "123",
			},
		},
	}

	config := &types.ResponseConfig{
		Mapping: map[string]string{
			"cnpj":       "$.data.cnpj",
			"razao":      "$.data.nome",
			"situacao":   "$.data.situacao",
			"logradouro": "$.data.endereco.logradouro",
		},
	}

	result, err := engine.Transform(jsonData, config)

	assert.NoError(t, err)
	assert.Equal(t, "12345678000195", result["cnpj"])
	assert.Equal(t, "Empresa Teste LTDA", result["razao"])
	assert.Equal(t, "01", result["situacao"])
	assert.Equal(t, "Rua Teste", result["logradouro"])
}

func TestEngine_Transform_WithTransformations(t *testing.T) {
	engine := NewEngine()
	engine.RegisterPlugin("format_cnpj", &FormatCNPJPlugin{})

	jsonData := map[string]interface{}{
		"data": map[string]interface{}{
			"cnpj":     "12345678000195",
			"situacao": "01",
		},
	}

	config := &types.ResponseConfig{
		Mapping: map[string]string{
			"cnpj":     "$.data.cnpj",
			"situacao": "$.data.situacao",
		},
		Transforms: []types.TransformConfig{
			{
				Field:     "cnpj",
				Operation: "format_cnpj",
			},
			{
				Field:     "situacao",
				Operation: "map_values",
				Values: map[string]string{
					"01": "ativa",
					"02": "suspensa",
				},
			},
		},
	}

	result, err := engine.Transform(jsonData, config)

	assert.NoError(t, err)
	assert.Equal(t, "12.345.678/0001-95", result["cnpj"])
	assert.Equal(t, "ativa", result["situacao"])
}

func TestEngine_Transform_JSONString(t *testing.T) {
	engine := NewEngine()

	jsonString := `{"data": {"id": "123", "name": "Test"}}`

	config := &types.ResponseConfig{
		Mapping: map[string]string{
			"id":   "$.data.id",
			"name": "$.data.name",
		},
	}

	result, err := engine.Transform(jsonString, config)

	assert.NoError(t, err)
	assert.Equal(t, "123", result["id"])
	assert.Equal(t, "Test", result["name"])
}

func TestFormatCNPJPlugin(t *testing.T) {
	plugin := &FormatCNPJPlugin{}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Valid CNPJ",
			input:    "12345678000195",
			expected: "12.345.678/0001-95",
		},
		{
			name:     "Already formatted",
			input:    "12.345.678/0001-95",
			expected: "12.345.678/0001-95",
		},
		{
			name:     "Invalid length",
			input:    "123456",
			expected: "123456", // Returns original
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.Transform(tt.input, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatCPFPlugin(t *testing.T) {
	plugin := &FormatCPFPlugin{}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Valid CPF",
			input:    "12345678901",
			expected: "123.456.789-01",
		},
		{
			name:     "Already formatted",
			input:    "123.456.789-01",
			expected: "123.456.789-01",
		},
		{
			name:     "Invalid length",
			input:    "12345",
			expected: "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.Transform(tt.input, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatCEPPlugin(t *testing.T) {
	plugin := &FormatCEPPlugin{}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Valid CEP",
			input:    "01310100",
			expected: "01310-100",
		},
		{
			name:     "Already formatted",
			input:    "01310-100",
			expected: "01310-100",
		},
		{
			name:     "Invalid length",
			input:    "123",
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.Transform(tt.input, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToUpperPlugin(t *testing.T) {
	plugin := &ToUpperPlugin{}

	result, err := plugin.Transform("joão silva", nil)
	assert.NoError(t, err)
	assert.Equal(t, "JOÃO SILVA", result)
}

func TestToLowerPlugin(t *testing.T) {
	plugin := &ToLowerPlugin{}

	result, err := plugin.Transform("JOÃO SILVA", nil)
	assert.NoError(t, err)
	assert.Equal(t, "joão silva", result)
}

func TestTrimPlugin(t *testing.T) {
	plugin := &TrimPlugin{}

	result, err := plugin.Transform("  teste  ", nil)
	assert.NoError(t, err)
	assert.Equal(t, "teste", result)
}

func TestEngine_MapValues(t *testing.T) {
	engine := NewEngine()

	jsonData := map[string]interface{}{
		"status": "A",
	}

	config := &types.ResponseConfig{
		Mapping: map[string]string{
			"status": "$.status",
		},
		Transforms: []types.TransformConfig{
			{
				Field:     "status",
				Operation: "map_values",
				Values: map[string]string{
					"A": "ativo",
					"I": "inativo",
				},
			},
		},
	}

	result, err := engine.Transform(jsonData, config)

	assert.NoError(t, err)
	assert.Equal(t, "ativo", result["status"])
}

func TestEngine_ExtractArray(t *testing.T) {
	engine := NewEngine()

	jsonData := map[string]interface{}{
		"data": map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"name": "Item 1"},
				map[string]interface{}{"name": "Item 2"},
			},
		},
	}

	config := &types.ResponseConfig{
		Mapping: map[string]string{
			"items": "$.data.items[*].name",
		},
	}

	result, err := engine.Transform(jsonData, config)

	assert.NoError(t, err)
	assert.IsType(t, []interface{}{}, result["items"])
	items := result["items"].([]interface{})
	assert.Len(t, items, 2)
	assert.Equal(t, "Item 1", items[0])
	assert.Equal(t, "Item 2", items[1])
}
