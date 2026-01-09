package middleware

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// setupTestRouter cria um router Gin para testes
func setupTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	return r
}

// TestNewFreemiumRateLimiter testa criação do rate limiter
func TestNewFreemiumRateLimiter(t *testing.T) {
	config := DefaultFreemiumConfig()
	limiter := NewFreemiumRateLimiter(nil, config)

	if limiter == nil {
		t.Error("NewFreemiumRateLimiter() should not return nil")
	}

	if limiter.config.FreeLimit != 5 {
		t.Errorf("FreeLimit = %d, expected 5", limiter.config.FreeLimit)
	}

	if limiter.cache == nil {
		t.Error("Cache should be initialized")
	}
}

// TestDefaultFreemiumConfig testa configuração padrão
func TestDefaultFreemiumConfig(t *testing.T) {
	config := DefaultFreemiumConfig()

	if config.FreeLimit != 5 {
		t.Errorf("FreeLimit = %d, expected 5", config.FreeLimit)
	}

	if !config.PremiumUnlimited {
		t.Error("PremiumUnlimited should be true")
	}

	if config.WindowDuration != 24*time.Hour {
		t.Errorf("WindowDuration = %v, expected 24h", config.WindowDuration)
	}
}

// TestFreemiumRateLimiter_FirstRequest testa primeira request
func TestFreemiumRateLimiter_FirstRequest(t *testing.T) {
	config := DefaultFreemiumConfig()
	limiter := NewFreemiumRateLimiter(nil, config)
	router := setupTestRouter(limiter.Middleware())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusOK)
	}

	// Verificar headers de rate limit
	if w.Header().Get("X-RateLimit-Limit") != "5" {
		t.Errorf("X-RateLimit-Limit = %s, expected 5", w.Header().Get("X-RateLimit-Limit"))
	}

	// Remaining mostra quantas requests restam APÓS esta request (count já foi incrementado)
	if w.Header().Get("X-RateLimit-Remaining") != "5" {
		t.Errorf("X-RateLimit-Remaining = %s, expected 5", w.Header().Get("X-RateLimit-Remaining"))
	}

	if w.Header().Get("X-RateLimit-Reset") == "" {
		t.Error("X-RateLimit-Reset should be present")
	}
}

// TestFreemiumRateLimiter_MultipleRequests testa múltiplas requests
func TestFreemiumRateLimiter_MultipleRequests(t *testing.T) {
	config := DefaultFreemiumConfig()
	limiter := NewFreemiumRateLimiter(nil, config)
	router := setupTestRouter(limiter.Middleware())

	// Fazer 5 requests (limite)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Status code = %d, expected %d", i+1, w.Code, http.StatusOK)
		}

		// Remaining mostra quantas requests restam ANTES de incrementar
		// Request 1: count=0, remaining=5; Request 2: count=1, remaining=4; etc.
		expectedRemaining := 5 - i
		expectedRemainingStr := string(rune('0' + expectedRemaining))
		if w.Header().Get("X-RateLimit-Remaining") != expectedRemainingStr {
			t.Errorf("Request %d: X-RateLimit-Remaining = %s, expected %s",
				i+1, w.Header().Get("X-RateLimit-Remaining"), expectedRemainingStr)
		}
	}

	// 6ª request deve ser bloqueada
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusTooManyRequests)
	}

	if w.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Errorf("X-RateLimit-Remaining = %s, expected 0", w.Header().Get("X-RateLimit-Remaining"))
	}
}

// TestFreemiumRateLimiter_RateLimitExceeded testa mensagem de erro
func TestFreemiumRateLimiter_RateLimitExceeded(t *testing.T) {
	config := FreemiumConfig{
		FreeLimit:        1, // Limite de 1 para facilitar teste
		PremiumUnlimited: true,
		WindowDuration:   24 * time.Hour,
	}
	limiter := NewFreemiumRateLimiter(nil, config)
	router := setupTestRouter(limiter.Middleware())

	// Primeira request - OK
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request status = %d, expected %d", w1.Code, http.StatusOK)
	}

	// Segunda request - deve ser bloqueada
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request status = %d, expected %d", w2.Code, http.StatusTooManyRequests)
	}

	// Verificar mensagem de erro
	if w2.Body.String() == "" {
		t.Error("Error response body should not be empty")
	}
}

// TestFreemiumRateLimiter_DifferentIPs testa IPs diferentes
func TestFreemiumRateLimiter_DifferentIPs(t *testing.T) {
	config := FreemiumConfig{
		FreeLimit:        1,
		PremiumUnlimited: true,
		WindowDuration:   24 * time.Hour,
	}
	limiter := NewFreemiumRateLimiter(nil, config)
	router := setupTestRouter(limiter.Middleware())

	// IP 1 - primeira request
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("IP1 first request status = %d, expected %d", w1.Code, http.StatusOK)
	}

	// IP 2 - primeira request (diferente, deve passar)
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("IP2 first request status = %d, expected %d", w2.Code, http.StatusOK)
	}

	// IP 1 - segunda request (deve ser bloqueada)
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("IP1 second request status = %d, expected %d", w3.Code, http.StatusTooManyRequests)
	}
}

// TestFreemiumRateLimiter_ResetAfterWindow testa reset após janela de tempo
func TestFreemiumRateLimiter_ResetAfterWindow(t *testing.T) {
	config := FreemiumConfig{
		FreeLimit:        1,
		PremiumUnlimited: true,
		WindowDuration:   50 * time.Millisecond, // Janela curta para teste
	}
	limiter := NewFreemiumRateLimiter(nil, config)
	router := setupTestRouter(limiter.Middleware())

	// Primeira request
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request status = %d, expected %d", w1.Code, http.StatusOK)
	}

	// Segunda request imediata (deve ser bloqueada)
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request status = %d, expected %d", w2.Code, http.StatusTooManyRequests)
	}

	// Aguardar janela expirar
	time.Sleep(60 * time.Millisecond)

	// Terceira request após reset (deve passar)
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Third request after reset status = %d, expected %d", w3.Code, http.StatusOK)
	}
}

// TestGetUserIdentifier_AuthenticatedUser testa identificação de usuário autenticado
func TestGetUserIdentifier_AuthenticatedUser(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "user123")
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	identifier := limiter.getUserIdentifier(c)

	expected := "user:user123"
	if identifier != expected {
		t.Errorf("getUserIdentifier() = %s, expected %s", identifier, expected)
	}
}

// TestGetUserIdentifier_AnonymousUser testa identificação por IP
func TestGetUserIdentifier_AnonymousUser(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.RemoteAddr = "192.168.1.1:12345"

	identifier := limiter.getUserIdentifier(c)

	expected := "ip:192.168.1.1"
	if identifier != expected {
		t.Errorf("getUserIdentifier() = %s, expected %s", identifier, expected)
	}
}

// TestGetClientIP_XForwardedFor testa extração de IP com X-Forwarded-For
func TestGetClientIP_XForwardedFor(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Forwarded-For", "203.0.113.1")

	ip := limiter.getClientIP(c)

	expected := "203.0.113.1"
	if ip != expected {
		t.Errorf("getClientIP() = %s, expected %s", ip, expected)
	}
}

// TestGetClientIP_XRealIP testa extração de IP com X-Real-IP
func TestGetClientIP_XRealIP(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Real-IP", "198.51.100.1")

	ip := limiter.getClientIP(c)

	expected := "198.51.100.1"
	if ip != expected {
		t.Errorf("getClientIP() = %s, expected %s", ip, expected)
	}
}

// TestGetClientIP_RemoteAddr testa fallback para RemoteAddr
func TestGetClientIP_RemoteAddr(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.RemoteAddr = "192.0.2.1:54321"

	ip := limiter.getClientIP(c)

	expected := "192.0.2.1"
	if ip != expected {
		t.Errorf("getClientIP() = %s, expected %s", ip, expected)
	}
}

// TestCheckPremiumStatus testa verificação de status premium
func TestCheckPremiumStatus(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// Por enquanto todos são free (implementação futura)
	isPremium := limiter.checkPremiumStatus(c, "user:123")

	if isPremium {
		t.Error("checkPremiumStatus() should return false (not implemented yet)")
	}
}

// TestCheckRateLimit_NewEntry testa verificação de nova entrada
func TestCheckRateLimit_NewEntry(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	allowed, remaining, resetAt := limiter.checkRateLimit("test-id")

	if !allowed {
		t.Error("checkRateLimit() should allow new entry")
	}

	if remaining != 5 {
		t.Errorf("Remaining = %d, expected 5", remaining)
	}

	if resetAt.IsZero() {
		t.Error("ResetAt should not be zero")
	}

	// Verificar que resetAt está no futuro
	if !resetAt.After(time.Now()) {
		t.Error("ResetAt should be in the future")
	}
}

// TestIncrementCounter testa incremento de contador
func TestIncrementCounter(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	identifier := "test-id"

	// Primeira verificação (cria entrada)
	_, remaining1, _ := limiter.checkRateLimit(identifier)
	if remaining1 != 5 {
		t.Errorf("Initial remaining = %d, expected 5", remaining1)
	}

	// Incrementar
	limiter.incrementCounter(identifier)

	// Verificar novamente
	_, remaining2, _ := limiter.checkRateLimit(identifier)
	if remaining2 != 4 {
		t.Errorf("After increment remaining = %d, expected 4", remaining2)
	}
}

// TestCleanupCache testa limpeza de cache
func TestCleanupCache(t *testing.T) {
	config := FreemiumConfig{
		FreeLimit:        5,
		PremiumUnlimited: true,
		WindowDuration:   10 * time.Millisecond, // Janela muito curta
	}
	limiter := NewFreemiumRateLimiter(nil, config)

	// Criar entrada
	limiter.checkRateLimit("test-id")

	// Verificar que entrada existe
	limiter.cacheMu.RLock()
	_, exists := limiter.cache["test-id"]
	limiter.cacheMu.RUnlock()

	if !exists {
		t.Error("Entry should exist in cache")
	}

	// Aguardar expiração
	time.Sleep(20 * time.Millisecond)

	// Forçar cleanup manualmente
	limiter.cacheMu.Lock()
	now := time.Now()
	for identifier, entry := range limiter.cache {
		if now.After(entry.resetAt) {
			delete(limiter.cache, identifier)
		}
	}
	limiter.cacheMu.Unlock()

	// Verificar que entrada foi removida
	limiter.cacheMu.RLock()
	_, exists = limiter.cache["test-id"]
	limiter.cacheMu.RUnlock()

	if exists {
		t.Error("Expired entry should be removed from cache")
	}
}

// TestFreemiumRateLimiter_CustomConfig testa configuração customizada
func TestFreemiumRateLimiter_CustomConfig(t *testing.T) {
	config := FreemiumConfig{
		FreeLimit:        10, // Limite customizado
		PremiumUnlimited: false,
		WindowDuration:   1 * time.Hour,
	}
	limiter := NewFreemiumRateLimiter(nil, config)

	if limiter.config.FreeLimit != 10 {
		t.Errorf("FreeLimit = %d, expected 10", limiter.config.FreeLimit)
	}

	if limiter.config.PremiumUnlimited {
		t.Error("PremiumUnlimited should be false")
	}

	if limiter.config.WindowDuration != 1*time.Hour {
		t.Errorf("WindowDuration = %v, expected 1h", limiter.config.WindowDuration)
	}
}

// TestFreemiumRateLimiter_ConcurrentAccess testa acesso concurrent
func TestFreemiumRateLimiter_ConcurrentAccess(t *testing.T) {
	config := DefaultFreemiumConfig()
	limiter := NewFreemiumRateLimiter(nil, config)
	router := setupTestRouter(limiter.Middleware())

	// Executar múltiplas goroutines concorrentes
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			done <- true
		}(i)
	}

	// Aguardar todas as goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Teste passou se não houver race condition (detectado por -race flag)
}

// TestFreemiumRateLimiter_ThreadSafety testa thread safety
func TestFreemiumRateLimiter_ThreadSafety(t *testing.T) {
	limiter := NewFreemiumRateLimiter(nil, DefaultFreemiumConfig())

	done := make(chan bool, 100)

	// Múltiplas goroutines acessando o cache simultaneamente
	for i := 0; i < 100; i++ {
		go func(id int) {
			limiter.checkRateLimit("test-id")
			limiter.incrementCounter("test-id")
			done <- true
		}(i)
	}

	// Aguardar todas as goroutines
	for i := 0; i < 100; i++ {
		<-done
	}

	// Teste passou se não houver race condition
}

// TestFreemiumRateLimiter_WithDatabase testa com database (mock)
func TestFreemiumRateLimiter_WithDatabase(t *testing.T) {
	// Mock database (nil é aceito pois não estamos usando ainda)
	var db *sql.DB = nil

	config := DefaultFreemiumConfig()
	limiter := NewFreemiumRateLimiter(db, config)

	if limiter == nil {
		t.Error("NewFreemiumRateLimiter() should handle nil database")
	}
}
