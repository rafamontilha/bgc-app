package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IdempotencyCache stores API responses for idempotency
type IdempotencyCache struct {
	mu    sync.RWMutex
	store map[string]*CachedResponse
}

// CachedResponse represents a cached API response
type CachedResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
	CachedAt   time.Time
	ExpiresAt  time.Time
}

// NewIdempotencyCache creates a new idempotency cache
func NewIdempotencyCache() *IdempotencyCache {
	cache := &IdempotencyCache{
		store: make(map[string]*CachedResponse),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// cleanupExpired removes expired entries periodically
func (ic *IdempotencyCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ic.mu.Lock()
		now := time.Now()
		for key, entry := range ic.store {
			if now.After(entry.ExpiresAt) {
				delete(ic.store, key)
			}
		}
		ic.mu.Unlock()
	}
}

// Get retrieves a cached response
func (ic *IdempotencyCache) Get(key string) (*CachedResponse, bool) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()

	entry, exists := ic.store[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry, true
}

// Set stores a response in the cache
func (ic *IdempotencyCache) Set(key string, response *CachedResponse) {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	ic.store[key] = response
}

// generateIdempotencyHash creates a hash from request details
func generateIdempotencyHash(method, path string, params map[string]string) string {
	// Create deterministic string from request
	data := fmt.Sprintf("%s:%s:%v", method, path, params)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// IdempotencyMiddleware provides idempotency for POST/PUT requests
type IdempotencyMiddleware struct {
	cache *IdempotencyCache
	ttl   time.Duration
}

// NewIdempotencyMiddleware creates a new idempotency middleware
func NewIdempotencyMiddleware(ttl time.Duration) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		cache: NewIdempotencyCache(),
		ttl:   ttl,
	}
}

// Handle processes idempotency for requests
func (im *IdempotencyMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to POST, PUT, PATCH methods
		if c.Request.Method != http.MethodPost &&
			c.Request.Method != http.MethodPut &&
			c.Request.Method != http.MethodPatch {
			c.Next()
			return
		}

		// Get idempotency key from header
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			// No idempotency key provided, proceed normally
			c.Next()
			return
		}

		// Validate idempotency key format (should be UUID or similar)
		if len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "INVALID_IDEMPOTENCY_KEY",
					"message":    "Idempotency-Key header must be between 16 and 128 characters",
					"request_id": c.GetString("request_id"),
				},
			})
			c.Abort()
			return
		}

		// Generate cache key combining idempotency key with request details
		queryParams := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				queryParams[key] = values[0]
			}
		}

		cacheKey := generateIdempotencyHash(
			c.Request.Method,
			c.Request.URL.Path,
			queryParams,
		) + ":" + idempotencyKey

		// Check if we have a cached response
		if cached, exists := im.cache.Get(cacheKey); exists {
			// Return cached response
			c.Header("X-Idempotency-Cached", "true")
			c.Header("X-Idempotency-Cached-At", cached.CachedAt.Format(time.RFC3339))

			// Apply cached headers
			for key, value := range cached.Headers {
				c.Header(key, value)
			}

			c.JSON(cached.StatusCode, cached.Body)
			c.Abort()
			return
		}

		// Store the idempotency key in context for handlers to use
		c.Set("idempotency_key", idempotencyKey)

		// Create a custom response writer to capture the response
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           &[]byte{},
		}
		c.Writer = blw

		// Process the request
		c.Next()

		// Cache the response if request was successful
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			var responseBody interface{}
			if len(*blw.body) > 0 {
				_ = json.Unmarshal(*blw.body, &responseBody)
			}

			cached := &CachedResponse{
				StatusCode: c.Writer.Status(),
				Body:       responseBody,
				Headers: map[string]string{
					"Content-Type": c.Writer.Header().Get("Content-Type"),
				},
				CachedAt:  time.Now(),
				ExpiresAt: time.Now().Add(im.ttl),
			}

			im.cache.Set(cacheKey, cached)
		}
	}
}

// bodyLogWriter captures response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *[]byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	*w.body = append(*w.body, b...)
	return w.ResponseWriter.Write(b)
}

// GetStats returns cache statistics
func (im *IdempotencyMiddleware) GetStats() map[string]interface{} {
	im.cache.mu.RLock()
	defer im.cache.mu.RUnlock()

	return map[string]interface{}{
		"total_cached":  len(im.cache.store),
		"cache_ttl_hrs": im.ttl.Hours(),
	}
}
