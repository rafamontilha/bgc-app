package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaValidator manages JSON schema validation
type SchemaValidator struct {
	schemas map[string]*gojsonschema.Schema
	mu      sync.RWMutex
}

// ValidationError represents a schema validation error
type ValidationError struct {
	Field   string `json:"field"`
	Issue   string `json:"issue"`
	Context string `json:"context,omitempty"`
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid  bool               `json:"valid"`
	Errors []ValidationError  `json:"errors,omitempty"`
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schemaDir string) (*SchemaValidator, error) {
	sv := &SchemaValidator{
		schemas: make(map[string]*gojsonschema.Schema),
	}

	// Load all schema files
	schemas := map[string]string{
		"market-size-request":      "market-size-request.schema.json",
		"market-size-response":     "market-size-response.schema.json",
		"route-comparison-request": "route-comparison-request.schema.json",
		"route-comparison-response": "route-comparison-response.schema.json",
		"error-response":           "error-response.schema.json",
	}

	for name, filename := range schemas {
		schemaPath := filepath.Join(schemaDir, filename)

		// Check if file exists
		if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("schema file not found: %s", schemaPath)
		}

		// Get absolute path for file:// URL
		absPath, err := filepath.Abs(schemaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for %s: %w", schemaPath, err)
		}

		schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", absPath))
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil {
			return nil, fmt.Errorf("failed to load schema %s: %w", name, err)
		}

		sv.schemas[name] = schema
	}

	return sv, nil
}

// ValidateRequest validates a request against a named schema
func (sv *SchemaValidator) ValidateRequest(schemaName string, data interface{}) (*ValidationResult, error) {
	sv.mu.RLock()
	schema, exists := sv.schemas[schemaName]
	sv.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("schema not found: %s", schemaName)
	}

	documentLoader := gojsonschema.NewGoLoader(data)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	valResult := &ValidationResult{
		Valid:  result.Valid(),
		Errors: []ValidationError{},
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			valResult.Errors = append(valResult.Errors, ValidationError{
				Field:   desc.Field(),
				Issue:   desc.Description(),
				Context: desc.Context().String(),
			})
		}
	}

	return valResult, nil
}

// ValidateMarketSizeRequest validates market size request parameters
func (sv *SchemaValidator) ValidateMarketSizeRequest(data interface{}) (*ValidationResult, error) {
	return sv.ValidateRequest("market-size-request", data)
}

// ValidateRouteComparisonRequest validates route comparison request parameters
func (sv *SchemaValidator) ValidateRouteComparisonRequest(data interface{}) (*ValidationResult, error) {
	return sv.ValidateRequest("route-comparison-request", data)
}

// GetAvailableSchemas returns list of loaded schema names
func (sv *SchemaValidator) GetAvailableSchemas() []string {
	sv.mu.RLock()
	defer sv.mu.RUnlock()

	names := make([]string, 0, len(sv.schemas))
	for name := range sv.schemas {
		names = append(names, name)
	}
	return names
}
