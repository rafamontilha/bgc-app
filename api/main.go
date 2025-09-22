package main

import (
	"fmt"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"  // remova se não estiver usando
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"   // <- blank import com ESPAÇO
	"gopkg.in/yaml.v3"
)

/* ===== Tipos de domínio ===== */

type MarketItem struct {
	Ano        int     `json:"ano"`
	NCMChapter string  `json:"ncm_chapter"`
	ValorUSD   float64 `json:"valor_usd"`
}

type AppConfig struct {
	ScopeChapters []string
	SOMBase       float64
	SOMAggressive float64
}

type PartnerWeights map[string]map[string]float64 // chapter -> (partner->share)

/* Tarifas (cenários) */
type TariffScenario struct {
	Default  map[string]float64            `yaml:"default"`
	Chapters map[string]map[string]float64 `yaml:"chapters"`
	Years    map[string]struct {
		Default  map[string]float64            `yaml:"default"`
		Chapters map[string]map[string]float64 `yaml:"chapters"`
	} `yaml:"years"`
}
type TariffScenarios struct {
	Scenarios map[string]TariffScenario `yaml:"scenarios"`
}

/* ===== Estado global ===== */

var (
	db       *sql.DB
	appCfg   AppConfig
	pweights PartnerWeights
	tariffs  TariffScenarios
)

/* ===== Utilidades ===== */

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustConnectDB() *sql.DB {
	host := getenv("DB_HOST", "db")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "bgc")
	pass := getenv("DB_PASS", "bgc")
	name := getenv("DB_NAME", "bgc")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	var conn *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		conn, err = sql.Open("postgres", dsn)
		if err == nil {
			if pingErr := conn.Ping(); pingErr == nil {
				log.Printf("Connected to Postgres at %s:%s", host, port)
				return conn
			} else {
				err = pingErr
			}
		}
		log.Printf("Waiting for Postgres... (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("Failed to connect to DB: %v", err)
	return nil
}

func loadConfig() AppConfig {
	cfg := AppConfig{
		ScopeChapters: []string{"02", "08", "84", "85"},
		SOMBase:       0.015,
		SOMAggressive: 0.03,
	}
	if v := getenv("SCOPE_CHAPTERS", ""); v != "" {
		parts := strings.Split(v, ",")
		clean := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if len(p) == 1 {
				p = "0" + p
			}
			if len(p) >= 2 {
				clean = append(clean, p[:2])
			}
		}
		if len(clean) > 0 {
			cfg.ScopeChapters = clean
		}
	}
	if v := getenv("SOM_BASE", ""); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.SOMBase = f
		}
	}
	if v := getenv("SOM_AGGRESSIVE", ""); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.SOMAggressive = f
		}
	}
	return cfg
}

func loadPartnerWeights() PartnerWeights {
	path := getenv("PARTNER_WEIGHTS_FILE", "./config/partners_stub.yaml")

	// Estrutura esperada no YAML: partners: { "<chapter|default>": { "<PAIS>": <peso float> } }
	var doc struct {
		Partners map[string]map[string]float64 `yaml:"partners"`
	}

	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("partner weights not found (%s): using defaults", path)
		return nil
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		log.Printf("failed to parse partner weights: %v", err)
		return nil
	}
	if len(doc.Partners) == 0 {
		log.Printf("partner weights file is empty: %s", path)
		return nil
	}

	out := make(PartnerWeights) // PartnerWeights == map[string]map[string]float64

	for chapterKey, partnersMap := range doc.Partners {
		// normaliza chave do capítulo (mantém "default" como está; capitulos "84" etc. sem mexer)
		chKey := strings.TrimSpace(chapterKey)
		if chKey != "default" && len(chKey) == 1 {
			chKey = "0" + chKey
		}

		// Cria o mapa interno para este capítulo caso ainda não exista
		if _, ok := out[chKey]; !ok {
			out[chKey] = make(map[string]float64)
		}

		for partnerCode, weight := range partnersMap {
			p := strings.ToUpper(strings.TrimSpace(partnerCode))
			out[chKey][p] = weight // <- aqui o tipo é float64, coerente com map[string]float64
		}
	}

	return out
}
func loadTariffs() {
	path := getenv("TARIFF_SCENARIOS_FILE", "./config/tariff_scenarios.yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("tariff scenarios not found (%s): continuing without tariffs", path)
		return
	}
	if err := yaml.Unmarshal(b, &tariffs); err != nil {
		log.Printf("failed to parse tariff scenarios: %v", err)
	}
}

/* Resolve fator de tarifa (prioridade: ano.capítulo → ano.default → capítulo → default → 1.0) */
func factorFor(scn TariffScenario, year int, chapter, partner string) float64 {
	p := strings.ToUpper(partner)
	chap := chapter
	ys := fmt.Sprintf("%d", year)

	if y, ok := scn.Years[ys]; ok {
		if mp, ok := y.Chapters[chap]; ok {
			if f, ok := mp[p]; ok {
				return f
			}
		}
		if f, ok := y.Default[p]; ok {
			return f
		}
	}
	if mp, ok := scn.Chapters[chap]; ok {
		if f, ok := mp[p]; ok {
			return f
		}
	}
	if f, ok := scn.Default[p]; ok {
		return f
	}
	return 1.0
}

/* ===== Handlers ===== */

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*" // quando a origem é vazia (ex.: curl), tudo bem
		} else {
			// Em dev, pode liberar geral. Se preferir restringir:
			// if origin == "http://localhost:3000" { ... }
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Não vamos usar cookies -> Credenciais desabilitadas
		// c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// ===== S2-16: Request ID, Logs JSON e Métricas =====

type routeMetrics struct {
	Requests     int64   `json:"requests"`
	Status2xx    int64   `json:"status_2xx"`
	Status4xx    int64   `json:"status_4xx"`
	Status5xx    int64   `json:"status_5xx"`
	SumLatencyMs int64   `json:"sum_latency_ms"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
}

var (
	metricsStart   = time.Now()
	metricsMu      sync.RWMutex
	totalRequests  int64
	statusCounters = map[int]int64{}              // ex.: 200->123, 404->7...
	byRoute        = map[string]*routeMetrics{}   // chave: "METHOD " + rota (FullPath)
)

func newReqID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b)
}

// X-Request-Id: usa o header do cliente se vier, senão gera um novo
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get("X-Request-Id")
		if rid == "" {
			rid = newReqID()
		}
		c.Set("req_id", rid)
		c.Writer.Header().Set("X-Request-Id", rid)
		c.Next()
	}
}

// Métricas + log estruturado (JSON)
func metricsAndLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latMs := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		method := c.Request.Method

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		key := method + " " + route

		metricsMu.Lock()
		totalRequests++
		statusCounters[status]++
		rm := byRoute[key]
		if rm == nil {
			rm = &routeMetrics{}
			byRoute[key] = rm
		}
		rm.Requests++
		rm.SumLatencyMs += latMs
		switch {
		case status >= 200 && status < 300:
			rm.Status2xx++
		case status >= 400 && status < 500:
			rm.Status4xx++
		case status >= 500:
			rm.Status5xx++
		}
		if rm.Requests > 0 {
			rm.AvgLatencyMs = float64(rm.SumLatencyMs) / float64(rm.Requests)
		}
		metricsMu.Unlock()

		// Log JSON minimalista
		rid, _ := c.Get("req_id")
		log.Printf(`{"ts":"%s","level":"info","req_id":"%v","method":"%s","path":"%s","route":"%s","status":%d,"latency_ms":%d,"ip":"%s","ua":"%s"}`,
			time.Now().Format(time.RFC3339Nano),
			rid, method, c.Request.URL.Path, route, status, latMs, c.ClientIP(), c.Request.UserAgent(),
		)
	}
}

// Handler do /metrics (snapshot em JSON)
func metricsHandler(c *gin.Context) {
	metricsMu.RLock()
	defer metricsMu.RUnlock()

	// copia defensiva para resposta
	statusByString := map[string]int64{}
	for code, cnt := range statusCounters {
		statusByString[strconv.Itoa(code)] = cnt
	}
	routesCopy := map[string]routeMetrics{}
	for k, v := range byRoute {
		routesCopy[k] = *v
	}

	c.JSON(http.StatusOK, gin.H{
		"uptime_seconds": int64(time.Since(metricsStart).Seconds()),
		"requests_total": totalRequests,
		"requests_by_status": statusByString,
		"routes": routesCopy,
	})
}

func healthHandler(c *gin.Context) {
    // opcional: checar DB rapidamente (descomente se quiser)
    // if err := db.Ping(); err != nil {
    //     c.JSON(500, gin.H{"status":"degraded","error": err.Error()})
    //     return
    // }
    c.JSON(200, gin.H{
        "select": 1,           // compat com a resposta antiga
        "status": "ok",
        "chapters_onda1":   appCfg.ScopeChapters,
        "partner_weights":  (pweights != nil),
        "tariffs_loaded":   (len(tariffs.Scenarios) > 0),
        "available_scenarios": func() []string {
            keys := make([]string, 0, len(tariffs.Scenarios))
            for k := range tariffs.Scenarios { keys = append(keys, k) }
            return keys
        }(),
    })
}


func main() {
	gin.SetMode(gin.ReleaseMode)
	appCfg = loadConfig()
	pweights = loadPartnerWeights()
	loadTariffs()

	db = mustConnectDB()
	defer db.Close()

	r := gin.Default()
	r.Use(corsMiddleware())
	r.Use(requestIDMiddleware())     // S2-16
	r.Use(metricsAndLogMiddleware()) // S2-16

	// --- ROTAS (sem duplicação) ---
	// Health (uma vez cada; ambos usam o mesmo handler)
	r.GET("/health",  healthHandler)
	r.GET("/healthz", healthHandler)

	// Métricas (S2-16)
	r.GET("/metrics", metricsHandler)

	// OpenAPI & ReDoc
	r.GET("/openapi.yaml", func(c *gin.Context) { c.File("./openapi.yaml") })
	r.GET("/docs", func(c *gin.Context) {
		html := `<!doctype html><html><head><meta charset="utf-8"><title>BGC API Docs</title></head>
<body><redoc spec-url='/openapi.yaml'></redoc>
<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script></body></html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// Endpoints analíticos
	r.GET("/market/size",    marketSizeHandler)
	r.GET("/routes/compare", routesCompareHandler)

	port := getenv("PORT", "8080")
	log.Printf("BGC API up on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}


func marketSizeHandler(c *gin.Context) {
	metric := strings.ToUpper(c.Query("metric"))
	if metric == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric is required [TAM|SAM|SOM]"})
		return
	}
	yearFrom, _ := strconv.Atoi(c.DefaultQuery("year_from", "2020"))
	yearTo, _ := strconv.Atoi(c.DefaultQuery("year_to", "2025"))
	ncmChapter := c.Query("ncm_chapter")
	scenario := strings.ToLower(c.DefaultQuery("scenario", "base"))

	var sb strings.Builder
	args := []any{yearFrom, yearTo}
	sb.WriteString(`SELECT ano, ncm_chapter, tam_total_usd
	                FROM v_tam_by_year_chapter
	               WHERE ano BETWEEN $1 AND $2`)
	argPos := 3

	if metric == "SAM" || metric == "SOM" {
		if len(appCfg.ScopeChapters) == 0 {
			c.JSON(500, gin.H{"error": "server misconfigured: scope chapters empty"})
			return
		}
		ph := make([]string, 0, len(appCfg.ScopeChapters))
		for range appCfg.ScopeChapters {
			ph = append(ph, fmt.Sprintf("$%d", argPos))
			argPos++
		}
		sb.WriteString(" AND ncm_chapter IN (" + strings.Join(ph, ",") + ")")
		for _, ch := range appCfg.ScopeChapters {
			args = append(args, ch)
		}
	}
	if ncmChapter != "" {
		sb.WriteString(fmt.Sprintf(" AND ncm_chapter = $%d", argPos))
		args = append(args, ncmChapter)
		argPos++
	}
	sb.WriteString(" ORDER BY ano, ncm_chapter")

	rows, err := db.Query(sb.String(), args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]MarketItem, 0, 64)
	for rows.Next() {
		var mi MarketItem
		var tam float64
		if err := rows.Scan(&mi.Ano, &mi.NCMChapter, &tam); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		switch metric {
		case "TAM", "SAM":
			mi.ValorUSD = tam
		case "SOM":
			switch scenario {
			case "aggressive":
				mi.ValorUSD = tam * appCfg.SOMAggressive
			default:
				mi.ValorUSD = tam * appCfg.SOMBase
			}
		default:
			c.JSON(400, gin.H{"error": "invalid metric; use TAM|SAM|SOM"})
			return
		}
		items = append(items, mi)
	}
	c.JSON(200, gin.H{"metric": metric, "scenario": scenario, "items": items})
}

func routesCompareHandler(c *gin.Context) {
	from := strings.ToUpper(c.DefaultQuery("from", "USA"))
	altsRaw := c.DefaultQuery("alts", "CHN,ARE,SAU,IND")
	alts := make([]string, 0)
	for _, a := range strings.Split(altsRaw, ",") {
		a = strings.TrimSpace(strings.ToUpper(a))
		if a != "" && a != from {
			alts = append(alts, a)
		}
	}
	year, err := strconv.Atoi(c.DefaultQuery("year", "2024"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid year"})
		return
	}
	chapter := c.Query("ncm_chapter")
	if len(chapter) == 1 {
		chapter = "0" + chapter
	}
	if chapter == "" || len(chapter) < 2 {
		c.JSON(400, gin.H{"error": "ncm_chapter (2 dígitos) é obrigatório"})
		return
	}

	// TAM base para ano/capítulo
	var tam float64
	q := `SELECT tam_total_usd FROM v_tam_by_year_chapter WHERE ano=$1 AND ncm_chapter=$2`
	if err := db.QueryRow(q, year, chapter).Scan(&tam); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "sem dados para ano/capítulo", "year": year, "ncm_chapter": chapter})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	partners := append([]string{from}, alts...)
	weights := map[string]float64{}

	// Pesos: capítulo → default → fallback 40/60
	if pweights != nil {
		if w, ok := pweights[chapter]; ok {
			for k, v := range w {
				weights[strings.ToUpper(k)] = v
			}
		}
		if len(weights) == 0 {
			if w, ok := pweights["default"]; ok {
				for k, v := range w {
					weights[strings.ToUpper(k)] = v
				}
			}
		}
	}
	if len(weights) == 0 {
		weights[from] = 0.40
		if len(alts) > 0 {
			rem := 0.60 / float64(len(alts))
			for _, a := range alts {
				weights[a] = rem
			}
		}
	}
	// normaliza ao conjunto solicitado
	sum := 0.0
	for _, p := range partners {
		sum += weights[p]
	}
	if sum == 0 {
		eq := 1.0 / float64(len(partners))
		for _, p := range partners {
			weights[p] = eq
		}
		sum = 1.0
	}
	for k, v := range weights {
		weights[k] = v / sum
	}

	// Tarifas (cenários)
	scenarioName := c.DefaultQuery("tariff_scenario", "base")
	scn, hasScenario := tariffs.Scenarios[scenarioName]
	tariffApplied := false

	type Item struct {
		Partner      string  `json:"partner"`
		Share        float64 `json:"share"`
		Factor       float64 `json:"factor"`
		EstimatedUSD float64 `json:"estimated_usd"`
	}

	out := make([]Item, 0, len(partners))
	adjustedTotal := 0.0
	for _, p := range partners {
		share := weights[p]
		est := share * tam

		factor := 1.0
		if hasScenario {
			factor = factorFor(scn, year, chapter, p)
			if factor != 1.0 {
				tariffApplied = true
			}
		}
		est = est * factor
		adjustedTotal += est

		out = append(out, Item{
			Partner:      p,
			Share:        share,
			Factor:       factor,
			EstimatedUSD: est,
		})
	}

	c.JSON(200, gin.H{
		"year":               year,
		"ncm_chapter":        chapter,
		"basis":              "TAM (mview)",
		"tam_total_usd":      tam,
		"from":               from,
		"alts":               alts,
		"tariff_scenario":    scenarioName,
		"tariff_applied":     tariffApplied,
		"adjusted_total_usd": adjustedTotal,
		"note":               "stub com pesos + fatores de tarifa; substituir por dados reais por parceiro em próxima onda",
		"results":            out,
	})
}
