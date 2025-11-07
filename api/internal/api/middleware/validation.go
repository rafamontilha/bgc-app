package middleware

import (
	"fmt"
	"net/http"

	"bgc-app/internal/api/validation"

	"github.com/gin-gonic/gin"
)

// ValidationMiddleware adds schema validation to routes
type ValidationMiddleware struct {
	validator *validation.SchemaValidator
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware(validator *validation.SchemaValidator) *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validator,
	}
}

// ValidateMarketSizeRequest validates market size request parameters
func (vm *ValidationMiddleware) ValidateMarketSizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract query parameters
		params := map[string]interface{}{
			"metric":      c.Query("metric"),
			"year_from":   c.Query("year_from"),
			"year_to":     c.Query("year_to"),
			"ncm_chapter": c.Query("ncm_chapter"),
			"scenario":    c.Query("scenario"),
		}

		// Remove empty optional parameters
		if params["ncm_chapter"] == "" {
			delete(params, "ncm_chapter")
		}
		if params["scenario"] == "" {
			delete(params, "scenario")
		}

		// Convert year_from and year_to to integers if present
		if yearFrom := c.Query("year_from"); yearFrom != "" {
			var year int
			if _, err := fmt.Sscanf(yearFrom, "%d", &year); err == nil {
				params["year_from"] = year
			}
		}
		if yearTo := c.Query("year_to"); yearTo != "" {
			var year int
			if _, err := fmt.Sscanf(yearTo, "%d", &year); err == nil {
				params["year_to"] = year
			}
		}

		// Validate against schema
		result, err := vm.validator.ValidateMarketSizeRequest(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":       "VALIDATION_INTERNAL_ERROR",
					"message":    "Failed to validate request",
					"request_id": c.GetString("request_id"),
				},
			})
			c.Abort()
			return
		}

		if !result.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "VALIDATION_ERROR",
					"message":    "Invalid request parameters",
					"details":    result.Errors,
					"request_id": c.GetString("request_id"),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateRouteComparisonRequest validates route comparison request parameters
func (vm *ValidationMiddleware) ValidateRouteComparisonRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract query parameters
		params := map[string]interface{}{
			"from":            c.Query("from"),
			"alternatives":    c.Query("alternatives"),
			"ncm_chapter":     c.Query("ncm_chapter"),
			"year":            c.Query("year"),
			"tariff_scenario": c.Query("tariff_scenario"),
		}

		// Remove empty optional parameters
		if params["tariff_scenario"] == "" {
			delete(params, "tariff_scenario")
		}

		// Convert year to integer if present
		if year := c.Query("year"); year != "" {
			var yearInt int
			if _, err := fmt.Sscanf(year, "%d", &yearInt); err == nil {
				params["year"] = yearInt
			}
		}

		// Validate against schema
		result, err := vm.validator.ValidateRouteComparisonRequest(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":       "VALIDATION_INTERNAL_ERROR",
					"message":    "Failed to validate request",
					"request_id": c.GetString("request_id"),
				},
			})
			c.Abort()
			return
		}

		if !result.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "VALIDATION_ERROR",
					"message":    "Invalid request parameters",
					"details":    result.Errors,
					"request_id": c.GetString("request_id"),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
