package app

import (
	"database/sql"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"bgc-app/internal/api/handlers"
	"bgc-app/internal/api/middleware"
	"bgc-app/internal/api/validation"
	"bgc-app/internal/business/destination"
	"bgc-app/internal/business/health"
	"bgc-app/internal/business/market"
	"bgc-app/internal/business/route"
	"bgc-app/internal/config"
	"bgc-app/internal/observability/metrics"
	"bgc-app/internal/repository/postgres"
)

type Server struct {
	router *gin.Engine
	config *config.AppConfig
}

func NewServer(cfg *config.AppConfig, db *sql.DB) *Server {
	gin.SetMode(gin.ReleaseMode)

	weights := config.LoadPartnerWeights(cfg.PartnerWeightsFile)
	tariffs := config.LoadTariffScenarios(cfg.TariffScenariosFile)

	marketRepo := postgres.NewMarketRepository(db)
	routeRepo := postgres.NewRouteRepository(db)
	destinationRepo := postgres.NewDestinationRepository(db)

	marketService := market.NewService(marketRepo, cfg)
	routeService := route.NewService(routeRepo, weights, tariffs)
	healthService := health.NewService(cfg, weights, tariffs)
	destinationService := destination.NewService(destinationRepo)

	marketHandler := handlers.NewMarketHandler(marketService)
	routeHandler := handlers.NewRouteHandler(routeService)
	healthHandler := handlers.NewHealthHandler(healthService)
	simulatorHandler := handlers.NewSimulatorHandler(destinationService)

	// Initialize schema validator
	schemaDir := filepath.Join(".", "schemas", "v1")
	validator, err := validation.NewSchemaValidator(schemaDir)
	if err != nil {
		log.Printf("Warning: Schema validator initialization failed: %v (validation disabled)", err)
		validator = nil
	} else {
		log.Printf("Schema validator initialized with schemas: %v", validator.GetAvailableSchemas())
	}

	// Initialize idempotency middleware with 24h TTL
	idempotencyMW := middleware.NewIdempotencyMiddleware(24 * time.Hour)
	log.Printf("Idempotency middleware initialized with 24h TTL")

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(otelgin.Middleware("bgc-api")) // OpenTelemetry tracing
	r.Use(middleware.MetricsAndLog())    // Structured logging
	r.Use(metrics.PrometheusMiddleware()) // Prometheus metrics

	// Start DB stats collector (updates every 15 seconds)
	metrics.StartDBStatsCollector(db, 15*time.Second)
	log.Printf("Database stats collector started (interval: 15s)")

	// Health and metrics endpoints (no versioning)
	r.GET("/health", healthHandler.GetHealth)
	r.GET("/healthz", healthHandler.GetHealth)

	// Prometheus metrics endpoint (native format)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Legacy JSON metrics endpoint (backwards compatibility)
	r.GET("/metrics/json", middleware.GetMetricsHandler())

	// Documentation endpoints
	r.GET("/openapi.yaml", func(c *gin.Context) { c.File("./openapi.yaml") })
	r.GET("/docs", func(c *gin.Context) {
		html := `<!doctype html><html><head><meta charset="utf-8"><title>BGC API Docs</title></head>
<body><redoc spec-url='/openapi.yaml'></redoc>
<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script></body></html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// API v1 group with idempotency support
	v1 := r.Group("/v1")
	v1.Use(idempotencyMW.Handle())

	// Apply validation middleware if validator is available
	if validator != nil {
		validationMW := middleware.NewValidationMiddleware(validator)
		v1.GET("/market/size", validationMW.ValidateMarketSizeRequest(), marketHandler.GetMarketSize)
		v1.GET("/routes/compare", validationMW.ValidateRouteComparisonRequest(), routeHandler.CompareRoutes)
	} else {
		// Fallback to routes without validation
		v1.GET("/market/size", marketHandler.GetMarketSize)
		v1.GET("/routes/compare", routeHandler.CompareRoutes)
	}

	// Simulator endpoints with freemium rate limiting
	freemiumMW := middleware.NewFreemiumRateLimiter(db, middleware.DefaultFreemiumConfig())
	log.Printf("Freemium rate limiter initialized (5 req/day for free tier)")

	simulator := v1.Group("/simulator")
	simulator.POST("/destinations", freemiumMW.Middleware(), simulatorHandler.SimulateDestinations)

	// Backwards compatibility: redirect legacy endpoints to v1
	r.GET("/market/size", func(c *gin.Context) {
		c.Redirect(301, "/v1/market/size?"+c.Request.URL.RawQuery)
	})
	r.GET("/routes/compare", func(c *gin.Context) {
		c.Redirect(301, "/v1/routes/compare?"+c.Request.URL.RawQuery)
	})

	return &Server{
		router: r,
		config: cfg,
	}
}

func (s *Server) Run() error {
	log.Printf("BGC API up on :%s", s.config.Port)
	return s.router.Run(":" + s.config.Port)
}
