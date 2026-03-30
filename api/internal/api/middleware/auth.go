package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwksKey representa uma chave pública no formato JWKS.
type jwksKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksDocument struct {
	Keys []jwksKey `json:"keys"`
}

// jwksCache armazena as chaves públicas do Clerk com TTL.
// Evita uma chamada HTTP por request — revalida apenas quando o cache expira.
type jwksCache struct {
	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
	ttl       time.Duration
	url       string
	client    *http.Client
}

func newJWKSCache(url string, ttl time.Duration) *jwksCache {
	return &jwksCache{
		keys:   make(map[string]*rsa.PublicKey),
		ttl:    ttl,
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// getKey retorna a chave pública para o kid informado.
// Atualiza o cache se estiver expirado.
func (c *jwksCache) getKey(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	if time.Since(c.fetchedAt) < c.ttl {
		if key, ok := c.keys[kid]; ok {
			c.mu.RUnlock()
			return key, nil
		}
	}
	c.mu.RUnlock()

	if err := c.refresh(); err != nil {
		return nil, fmt.Errorf("jwks refresh: %w", err)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key, ok := c.keys[kid]
	if !ok {
		return nil, fmt.Errorf("kid %q not found in JWKS", kid)
	}
	return key, nil
}

func (c *jwksCache) refresh() error {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return fmt.Errorf("fetch JWKS from %s: %w", c.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned HTTP %d", resp.StatusCode)
	}

	var doc jwksDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return fmt.Errorf("decode JWKS: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(doc.Keys))
	for _, k := range doc.Keys {
		if k.Kty != "RSA" || k.Use != "sig" {
			continue
		}
		pub, err := rsaPublicKeyFromJWK(k.N, k.E)
		if err != nil {
			continue
		}
		keys[k.Kid] = pub
	}

	c.mu.Lock()
	c.keys = keys
	c.fetchedAt = time.Now()
	c.mu.Unlock()

	return nil
}

// rsaPublicKeyFromJWK constrói uma *rsa.PublicKey a partir dos campos n e e do JWKS.
func rsaPublicKeyFromJWK(nB64, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, fmt.Errorf("decode n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, fmt.Errorf("decode e: %w", err)
	}

	// Converte bytes big-endian para int
	e := 0
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}

// ClerkAuth é o middleware de autenticação JWT para o Clerk.
type ClerkAuth struct {
	cache *jwksCache
}

// NewClerkAuth cria o middleware com cache de chaves JWKS renovado a cada hora.
func NewClerkAuth(jwksURL string) *ClerkAuth {
	return &ClerkAuth{
		cache: newJWKSCache(jwksURL, 1*time.Hour),
	}
}

// newClerkAuthWithCache permite injetar um cache customizado (usado em testes).
func newClerkAuthWithCache(cache *jwksCache) *ClerkAuth {
	return &ClerkAuth{cache: cache}
}

// Middleware retorna o handler Gin que valida tokens JWT do Clerk.
// Endpoints protegidos retornam 401 caso o token esteja ausente, inválido ou expirado.
// Em caso de sucesso, seta "user_id" (sub) e "user_email" no contexto Gin.
func (a *ClerkAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := extractBearer(c.GetHeader("Authorization"))
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authorization header with Bearer token required",
			})
			return
		}

		claims, err := a.validate(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			return
		}

		if sub, _ := claims["sub"].(string); sub != "" {
			c.Set("user_id", sub)
		}
		// email pode estar no claim padrão do Clerk
		if email, _ := claims["email"].(string); email != "" {
			c.Set("user_email", email)
		}

		c.Next()
	}
}

// OptionalMiddleware tenta extrair o usuário do JWT mas não rejeita a request
// caso o token esteja ausente. Usado no endpoint do simulador (freemium público).
func (a *ClerkAuth) OptionalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := extractBearer(c.GetHeader("Authorization"))
		if raw == "" {
			c.Next()
			return
		}

		claims, err := a.validate(raw)
		if err != nil {
			c.Next()
			return
		}

		if sub, _ := claims["sub"].(string); sub != "" {
			c.Set("user_id", sub)
		}
		if email, _ := claims["email"].(string); email != "" {
			c.Set("user_email", email)
		}

		c.Next()
	}
}

func (a *ClerkAuth) validate(tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		kid, ok := t.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("missing kid in token header")
		}
		return a.cache.getKey(kid)
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	return claims, nil
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
