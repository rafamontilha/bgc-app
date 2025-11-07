# API Schemas v1

JSON Schema definitions for BGC App API v1 endpoints.

## Overview

This directory contains versioned JSON schemas for request/response validation and documentation.

## Schemas

### Market Size API
- `market-size-request.schema.json` - Request parameters for `/v1/market/size`
- `market-size-response.schema.json` - Response structure for market size calculations

### Route Comparison API
- `route-comparison-request.schema.json` - Request parameters for `/v1/routes/compare`
- `route-comparison-response.schema.json` - Response structure for route comparisons

### Common
- `error-response.schema.json` - Standard error response format

## Usage

### Validation in Go API

```go
import "github.com/xeipuuv/gojsonschema"

// Load schema
schemaLoader := gojsonschema.NewReferenceLoader("file://./schemas/v1/market-size-request.schema.json")
schema, err := gojsonschema.NewSchema(schemaLoader)

// Validate request
documentLoader := gojsonschema.NewGoLoader(requestData)
result, err := schema.Validate(documentLoader)

if !result.Valid() {
    // Handle validation errors
    for _, err := range result.Errors() {
        fmt.Printf("- %s: %s\n", err.Field(), err.Description())
    }
}
```

### Validation in TypeScript/JavaScript

```typescript
import Ajv from 'ajv';
import marketSizeRequestSchema from './schemas/v1/market-size-request.schema.json';

const ajv = new Ajv();
const validate = ajv.compile(marketSizeRequestSchema);

const isValid = validate(requestData);
if (!isValid) {
  console.log(validate.errors);
}
```

## Schema Versioning

- **v1**: Initial schema version
- Future versions (v2, v3, etc.) will be added in separate directories
- Breaking changes require a new version
- Non-breaking changes can be added to existing version

## Validation Rules

### Market Size Request
- `metric`: Required, must be "TAM", "SAM", or "SOM"
- `year_from`: Required, integer between 2000-2100
- `year_to`: Required, integer between 2000-2100, must be >= year_from
- `ncm_chapter`: Optional, 2-digit string (e.g., "84")
- `scenario`: Optional, "base" or "aggressive" (default: "base")

### Route Comparison Request
- `from`: Required, 3-letter country code (ISO 3166-1 alpha-3)
- `alternatives`: Required, comma-separated list of 3-letter codes
- `ncm_chapter`: Required, 2-digit string
- `year`: Required, integer between 2000-2100
- `tariff_scenario`: Optional, "base", "tarifa10", or "tarifa20"

## Examples

### Market Size Request
```json
{
  "metric": "SOM",
  "year_from": 2020,
  "year_to": 2023,
  "ncm_chapter": "84",
  "scenario": "aggressive"
}
```

### Route Comparison Request
```json
{
  "from": "USA",
  "alternatives": "CHN,ARE,DEU",
  "ncm_chapter": "85",
  "year": 2023,
  "tariff_scenario": "base"
}
```

## Contributing

When adding new schemas:
1. Follow JSON Schema draft-07 specification
2. Include `$schema`, `$id`, `title`, and `description`
3. Add examples in the schema
4. Document in this README
5. Update API code to validate against new schemas
