package app

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	"bgc-app/internal/api/handlers"
	"bgc-app/internal/api/middleware"
	"bgc-app/internal/business/health"
	"bgc-app/internal/business/market"
	"bgc-app/internal/business/route"
	"bgc-app/internal/config"
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

	marketService := market.NewService(marketRepo, cfg)
	routeService := route.NewService(routeRepo, weights, tariffs)
	healthService := health.NewService(cfg, weights, tariffs)

	marketHandler := handlers.NewMarketHandler(marketService)
	routeHandler := handlers.NewRouteHandler(routeService)
	healthHandler := handlers.NewHealthHandler(healthService)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.MetricsAndLog())

	r.GET("/health", healthHandler.GetHealth)
	r.GET("/healthz", healthHandler.GetHealth)

	r.GET("/metrics", middleware.GetMetricsHandler())

	r.GET("/openapi.yaml", func(c *gin.Context) { c.File("./openapi.yaml") })
	r.GET("/docs", func(c *gin.Context) {
		html := `<!doctype html><html><head><meta charset="utf-8"><title>BGC API Docs</title></head>
<body><redoc spec-url='/openapi.yaml'></redoc>
<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script></body></html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/market/size", marketHandler.GetMarketSize)
	r.GET("/routes/compare", routeHandler.CompareRoutes)

	return &Server{
		router: r,
		config: cfg,
	}
}

func (s *Server) Run() error {
	log.Printf("BGC API up on :%s", s.config.Port)
	return s.router.Run(":" + s.config.Port)
}
