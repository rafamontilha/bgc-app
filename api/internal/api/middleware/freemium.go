package middleware

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// FreemiumConfig configuração do middleware freemium
type FreemiumConfig struct {
	FreeLimit        int           // Limite para usuários free (default: 5/dia)
	PremiumUnlimited bool          // Premium tem limite ilimitado (default: true)
	WindowDuration   time.Duration // Janela de tempo (default: 24h)
}

// DefaultFreemiumConfig retorna configuração padrão
func DefaultFreemiumConfig() FreemiumConfig {
	return FreemiumConfig{
		FreeLimit:        5,
		PremiumUnlimited: true,
		WindowDuration:   24 * time.Hour,
	}
}

// FreemiumRateLimiter middleware de rate limiting freemium
type FreemiumRateLimiter struct {
	config FreemiumConfig
	db     *sql.DB

	// Cache in-memory para reduzir hits no banco
	cache   map[string]*rateLimitEntry
	cacheMu sync.RWMutex
}

type rateLimitEntry struct {
	count     int
	resetAt   time.Time
	isPremium bool
}

// NewFreemiumRateLimiter cria novo middleware
func NewFreemiumRateLimiter(db *sql.DB, config FreemiumConfig) *FreemiumRateLimiter {
	limiter := &FreemiumRateLimiter{
		config: config,
		db:     db,
		cache:  make(map[string]*rateLimitEntry),
	}

	// Goroutine para limpeza de cache
	go limiter.cleanupCache()

	return limiter
}

// Middleware retorna o middleware Gin
func (f *FreemiumRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Identifica usuário (por IP ou user_id)
		identifier := f.getUserIdentifier(c)

		// Verifica se é premium
		isPremium := f.checkPremiumStatus(c, identifier)

		// Premium não tem limite
		if isPremium {
			c.Header("X-RateLimit-Limit", "unlimited")
			c.Header("X-RateLimit-Remaining", "unlimited")
			c.Next()
			return
		}

		// Verifica rate limit para free
		allowed, remaining, resetAt := f.checkRateLimit(identifier)

		// Headers de rate limit
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", f.config.FreeLimit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt.Unix()))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": fmt.Sprintf("Limite de %d simulações por dia atingido. Faça upgrade para Pro para simulações ilimitadas.", f.config.FreeLimit),
				"remaining": 0,
				"reset_at": resetAt.Unix(),
			})
			c.Abort()
			return
		}

		// Incrementa contador
		f.incrementCounter(identifier)

		c.Next()
	}
}

// getUserIdentifier identifica o usuário (user_id ou IP)
func (f *FreemiumRateLimiter) getUserIdentifier(c *gin.Context) string {
	// Tenta pegar user_id do context (se autenticado)
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%s", userID)
	}

	// Fallback para IP address
	ip := f.getClientIP(c)
	return fmt.Sprintf("ip:%s", ip)
}

// getClientIP pega o IP real do cliente (considerando proxies)
func (f *FreemiumRateLimiter) getClientIP(c *gin.Context) string {
	// X-Forwarded-For (se atrás de proxy)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// Pega o primeiro IP (cliente real)
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}

	// X-Real-IP
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// RemoteAddr (direto)
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}

	return c.Request.RemoteAddr
}

// checkPremiumStatus verifica se usuário é premium
func (f *FreemiumRateLimiter) checkPremiumStatus(c *gin.Context, identifier string) bool {
	// TODO: Implementar verificação real de premium
	// Por enquanto, todos são free
	//
	// Exemplo de implementação:
	// if userID, exists := c.Get("user_id"); exists {
	//     query := "SELECT is_premium FROM users WHERE id = $1"
	//     var isPremium bool
	//     f.db.QueryRow(query, userID).Scan(&isPremium)
	//     return isPremium
	// }
	return false
}

// checkRateLimit verifica se usuário está dentro do limite
func (f *FreemiumRateLimiter) checkRateLimit(identifier string) (allowed bool, remaining int, resetAt time.Time) {
	f.cacheMu.RLock()
	entry, exists := f.cache[identifier]
	f.cacheMu.RUnlock()

	now := time.Now()

	// Se não existe ou expirou, cria nova entrada
	if !exists || now.After(entry.resetAt) {
		f.cacheMu.Lock()
		f.cache[identifier] = &rateLimitEntry{
			count:   0,
			resetAt: now.Add(f.config.WindowDuration),
		}
		entry = f.cache[identifier]
		f.cacheMu.Unlock()
	}

	// Verifica limite
	remaining = f.config.FreeLimit - entry.count
	if remaining <= 0 {
		return false, 0, entry.resetAt
	}

	return true, remaining, entry.resetAt
}

// incrementCounter incrementa o contador de uso
func (f *FreemiumRateLimiter) incrementCounter(identifier string) {
	f.cacheMu.Lock()
	defer f.cacheMu.Unlock()

	if entry, exists := f.cache[identifier]; exists {
		entry.count++
	}

	// TODO: Persistir no banco para análise
	// INSERT INTO simulator_usage (identifier, timestamp)
	// VALUES ($1, now())
}

// cleanupCache limpa entradas expiradas do cache
func (f *FreemiumRateLimiter) cleanupCache() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		f.cacheMu.Lock()
		now := time.Now()
		for identifier, entry := range f.cache {
			if now.After(entry.resetAt) {
				delete(f.cache, identifier)
			}
		}
		f.cacheMu.Unlock()
	}
}
