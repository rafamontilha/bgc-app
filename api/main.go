package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

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

var db *sql.DB
var appCfg AppConfig
var pweights PartnerWeights

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
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("partner weights not found (%s): using defaults", path)
		return nil
	}
	var doc struct {
		Partners map[string]map[string]float64 `yaml:"partners"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		log.Printf("failed to parse partner weights: %v", err)
		return nil
	}
	out := PartnerWeights{}
	for chapter, mp := range doc.Partners {
		out[chapter] = map[string]float64{}
		for k, v := range mp {
			out[chapter][strings.ToUpper(strings.TrimSpace(k))] = v
		}
	}
	return out
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	appCfg = loadConfig()
	pweights = loadPartnerWeights()
	db = mustConnectDB()
	defer db.Close()

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok":              true,
			"chapters_onda1":  appCfg.ScopeChapters,
			"partner_weights": (pweights != nil),
		})
	})

	r.GET("/openapi.yaml", func(c *gin.Context) { c.File("./openapi.yaml") })
	r.GET("/docs", func(c *gin.Context) {
		html := `<!doctype html>
<html>
<head><meta charset="utf-8"><title>BGC API Docs</title></head>
<body>
  <redoc spec-url='/openapi.yaml'></redoc>
  <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/market/size", marketSizeHandler)
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

	// 1) tenta pesos específicos do capítulo; 2) default; 3) fallback 40/60
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

	// normaliza para o conjunto solicitado
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

	type Item struct {
		Partner      string  `json:"partner"`
		Share        float64 `json:"share"`
		EstimatedUSD float64 `json:"estimated_usd"`
	}
	out := make([]Item, 0, len(partners))
	for _, p := range partners {
		sh := weights[p]
		out = append(out, Item{
			Partner:      p,
			Share:        sh,
			EstimatedUSD: sh * tam,
		})
	}
	c.JSON(200, gin.H{
		"year":          year,
		"ncm_chapter":   chapter,
		"basis":         "TAM (mview)",
		"tam_total_usd": tam,
		"from":          from,
		"alts":          alts,
		"note":          "stub usando pesos configuráveis; trocar por dado real por parceiro na próxima iteração",
		"results":       out,
	})
}
