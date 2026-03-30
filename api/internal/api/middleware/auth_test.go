package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- helpers de teste ---

func newTestRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

// mintToken cria um JWT assinado com a chave informada.
func mintToken(t *testing.T, key *rsa.PrivateKey, kid, sub string, exp time.Time) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
		"iss": "https://clerk.test",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = kid

	signed, err := tok.SignedString(key)
	require.NoError(t, err)
	return signed
}

// jwksServerFor levanta um httptest.Server que serve o JWKS da chave RSA informada.
func jwksServerFor(t *testing.T, kid string, pub *rsa.PublicKey) *httptest.Server {
	t.Helper()

	nBytes := pub.N.Bytes()
	eBig := big.NewInt(int64(pub.E))
	eBytes := eBig.Bytes()

	jwksBody := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"kid": kid,
				"n":   base64.RawURLEncoding.EncodeToString(nBytes),
				"e":   base64.RawURLEncoding.EncodeToString(eBytes),
			},
		},
	}

	body, err := json.Marshal(jwksBody)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// ginContext cria um contexto Gin de teste com o header Authorization setado.
func ginContext(t *testing.T, authHeader string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		c.Request.Header.Set("Authorization", authHeader)
	}
	return c, w
}

// --- testes de rsaPublicKeyFromJWK ---

func TestRSAPublicKeyFromJWK_RoundTrip(t *testing.T) {
	key := newTestRSAKey(t)
	pub := &key.PublicKey

	nEnc := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	eBig := big.NewInt(int64(pub.E))
	eEnc := base64.RawURLEncoding.EncodeToString(eBig.Bytes())

	got, err := rsaPublicKeyFromJWK(nEnc, eEnc)
	require.NoError(t, err)
	assert.Equal(t, pub.N, got.N)
	assert.Equal(t, pub.E, got.E)
}

func TestRSAPublicKeyFromJWK_InvalidN(t *testing.T) {
	_, err := rsaPublicKeyFromJWK("!invalid base64", "AQAB")
	assert.Error(t, err)
}

func TestRSAPublicKeyFromJWK_InvalidE(t *testing.T) {
	key := newTestRSAKey(t)
	nEnc := base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes())
	_, err := rsaPublicKeyFromJWK(nEnc, "!invalid base64")
	assert.Error(t, err)
}

// --- testes de extractBearer ---

func TestExtractBearer(t *testing.T) {
	tests := []struct {
		header string
		want   string
	}{
		{"Bearer abc123", "abc123"},
		{"bearer abc123", "abc123"},
		{"BEARER abc123", "abc123"},
		{"Bearer ", ""},
		{"Token abc123", ""},
		{"", ""},
		{"abc123", ""},
	}

	for _, tc := range tests {
		got := extractBearer(tc.header)
		assert.Equal(t, tc.want, got, "header: %q", tc.header)
	}
}

// --- testes do ClerkAuth Middleware ---

func TestClerkAuth_ValidToken(t *testing.T) {
	key := newTestRSAKey(t)
	const kid = "test-key-1"

	jwksSrv := jwksServerFor(t, kid, &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	token := mintToken(t, key, kid, "user_clerk_abc", time.Now().Add(1*time.Hour))

	c, w := ginContext(t, "Bearer "+token)
	next := false
	auth.Middleware()(c)
	// gin.Context.Abort seta o index; se não abortou, next seria chamado
	// Verificamos diretamente o status do writer e o context
	if !c.IsAborted() {
		next = true
	}

	assert.True(t, next, "handler deve prosseguir com token válido")
	assert.Equal(t, 200, w.Code)

	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, "user_clerk_abc", userID)
}

func TestClerkAuth_MissingToken(t *testing.T) {
	key := newTestRSAKey(t)
	jwksSrv := jwksServerFor(t, "kid1", &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	c, w := ginContext(t, "")
	auth.Middleware()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestClerkAuth_ExpiredToken(t *testing.T) {
	key := newTestRSAKey(t)
	const kid = "kid-exp"
	jwksSrv := jwksServerFor(t, kid, &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	token := mintToken(t, key, kid, "user_x", time.Now().Add(-1*time.Hour)) // expirado

	c, w := ginContext(t, "Bearer "+token)
	auth.Middleware()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestClerkAuth_WrongKey(t *testing.T) {
	signingKey := newTestRSAKey(t)  // chave que assinou
	differentKey := newTestRSAKey(t) // chave errada no JWKS

	const kid = "kid-wrong"
	jwksSrv := jwksServerFor(t, kid, &differentKey.PublicKey) // publica errada
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	token := mintToken(t, signingKey, kid, "user_x", time.Now().Add(1*time.Hour))

	c, w := ginContext(t, "Bearer "+token)
	auth.Middleware()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestClerkAuth_UnknownKID(t *testing.T) {
	key := newTestRSAKey(t)
	wrongKey := newTestRSAKey(t)

	jwksSrv := jwksServerFor(t, "registered-kid", &wrongKey.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	token := mintToken(t, key, "unknown-kid", "user_x", time.Now().Add(1*time.Hour))

	c, w := ginContext(t, "Bearer "+token)
	auth.Middleware()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestClerkAuth_MalformedToken(t *testing.T) {
	key := newTestRSAKey(t)
	jwksSrv := jwksServerFor(t, "kid1", &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	c, w := ginContext(t, "Bearer notavalidjwt")
	auth.Middleware()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- testes do OptionalMiddleware ---

func TestClerkAuth_Optional_NoToken_Passes(t *testing.T) {
	key := newTestRSAKey(t)
	jwksSrv := jwksServerFor(t, "kid1", &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	c, w := ginContext(t, "")
	auth.OptionalMiddleware()(c)

	assert.False(t, c.IsAborted(), "sem token deve prosseguir normalmente")
	assert.Equal(t, 200, w.Code)
	_, exists := c.Get("user_id")
	assert.False(t, exists, "user_id não deve ser setado sem token")
}

func TestClerkAuth_Optional_ValidToken_SetsUserID(t *testing.T) {
	key := newTestRSAKey(t)
	const kid = "opt-kid"
	jwksSrv := jwksServerFor(t, kid, &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	token := mintToken(t, key, kid, "user_optional_123", time.Now().Add(1*time.Hour))

	c, _ := ginContext(t, "Bearer "+token)
	auth.OptionalMiddleware()(c)

	assert.False(t, c.IsAborted())
	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, "user_optional_123", userID)
}

func TestClerkAuth_Optional_InvalidToken_Passes(t *testing.T) {
	key := newTestRSAKey(t)
	jwksSrv := jwksServerFor(t, "kid1", &key.PublicKey)
	cache := newJWKSCache(jwksSrv.URL, 1*time.Hour)
	auth := newClerkAuthWithCache(cache)

	c, w := ginContext(t, "Bearer invalidtoken")
	auth.OptionalMiddleware()(c)

	assert.False(t, c.IsAborted(), "token inválido no optional não deve rejeitar")
	assert.Equal(t, 200, w.Code)
}

// --- teste do JWKS cache ---

func TestJWKSCache_RefreshOnExpiry(t *testing.T) {
	key := newTestRSAKey(t)
	const kid = "cache-kid"

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		nBytes := key.PublicKey.N.Bytes()
		eBytes := big.NewInt(int64(key.PublicKey.E)).Bytes()
		body, _ := json.Marshal(map[string]interface{}{
			"keys": []map[string]interface{}{{
				"kty": "RSA", "use": "sig", "alg": "RS256", "kid": kid,
				"n": base64.RawURLEncoding.EncodeToString(nBytes),
				"e": base64.RawURLEncoding.EncodeToString(eBytes),
			}},
		})
		w.Write(body)
	}))
	t.Cleanup(srv.Close)

	cache := newJWKSCache(srv.URL, 50*time.Millisecond) // TTL curto para o teste

	// Primeira chamada: deve buscar no servidor
	_, err := cache.getKey(kid)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Segunda chamada imediata: deve usar cache
	_, err = cache.getKey(kid)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount, "deve usar cache sem nova chamada HTTP")

	// Após expirar o TTL: deve buscar novamente
	time.Sleep(60 * time.Millisecond)
	_, err = cache.getKey(kid)
	require.NoError(t, err)
	assert.Equal(t, 2, callCount, "deve renovar após TTL expirar")
}
