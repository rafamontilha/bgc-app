package main

import (
	"log"
	"os"

	"github.com/bgc/integration-gateway/internal/auth"
	"github.com/bgc/integration-gateway/internal/framework"
	"github.com/bgc/integration-gateway/internal/observability"
	"github.com/bgc/integration-gateway/internal/registry"
	"github.com/bgc/integration-gateway/internal/transform"
	"github.com/bgc/integration-gateway/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Configuração
	configDir := getEnv("CONFIG_DIR", "./config/connectors")
	certsDir := getEnv("CERTS_DIR", "./certs")
	port := getEnv("PORT", "8081")
	environment := getEnv("ENVIRONMENT", "development")
	logLevel := getEnv("LOG_LEVEL", "info")

	// Configura log level
	observability.SetLogLevel(logLevel)

	observability.Info("Starting Integration Gateway",
		"config_dir", configDir,
		"certs_dir", certsDir,
		"environment", environment,
		"log_level", logLevel,
	)

	// Inicializa componentes
	reg := registry.NewRegistry(configDir)
	if err := reg.LoadAll(); err != nil {
		observability.Error("Failed to load connectors", "error", err)
		log.Fatalf("Failed to load connectors: %v", err)
	}
	observability.Info("Connectors loaded successfully", "count", reg.Count())

	certManager := auth.NewSimpleCertificateManager(certsDir)
	secretStore := auth.NewSimpleSecretStore()
	authEngine := auth.NewEngine(certManager, secretStore)

	transformEngine := transform.NewEngine()
	// Registra built-in plugins
	transformEngine.RegisterPlugin("format_cnpj", &transform.FormatCNPJPlugin{})
	transformEngine.RegisterPlugin("format_cpf", &transform.FormatCPFPlugin{})
	transformEngine.RegisterPlugin("format_cep", &transform.FormatCEPPlugin{})
	transformEngine.RegisterPlugin("to_upper", &transform.ToUpperPlugin{})
	transformEngine.RegisterPlugin("to_lower", &transform.ToLowerPlugin{})
	transformEngine.RegisterPlugin("trim", &transform.TrimPlugin{})

	executor := framework.NewExecutor(reg, authEngine, transformEngine)

	// Configura Gin
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":     "healthy",
			"connectors": reg.Count(),
		})
	})

	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Lista conectores
	router.GET("/v1/connectors", func(c *gin.Context) {
		connectors := reg.List()
		result := make([]gin.H, 0, len(connectors))
		for _, conn := range connectors {
			result = append(result, gin.H{
				"id":       conn.ID,
				"name":     conn.Name,
				"version":  conn.Version,
				"provider": conn.Provider,
				"endpoints": getEndpointNames(conn),
			})
		}
		c.JSON(200, result)
	})

	// Detalhes de um connector
	router.GET("/v1/connectors/:id", func(c *gin.Context) {
		id := c.Param("id")
		conn, err := reg.Get(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "connector not found"})
			return
		}
		c.JSON(200, conn)
	})

	// Executa endpoint de um connector
	router.POST("/v1/connectors/:id/:endpoint", func(c *gin.Context) {
		connectorID := c.Param("id")
		endpointName := c.Param("endpoint")

		// Parse request body (params)
		var params map[string]interface{}
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		// Executa
		ctx := &types.ExecutionContext{
			ConnectorID:  connectorID,
			EndpointName: endpointName,
			Environment:  environment,
			Params:       params,
		}

		result, err := executor.Execute(ctx)
		if err != nil {
			errorResponse := gin.H{"error": err.Error()}
			if result != nil {
				errorResponse["duration"] = result.Duration.String()
			}
			c.JSON(500, errorResponse)
			return
		}

		c.JSON(200, gin.H{
			"data":        result.Data,
			"status_code": result.StatusCode,
			"duration":    result.Duration.String(),
		})
	})

	// Inicia servidor
	addr := ":" + port
	observability.Info("Server starting", "address", addr)
	log.Printf("Server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		observability.Error("Failed to start server", "error", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEndpointNames(conn *types.ConnectorConfig) []string {
	names := make([]string, 0, len(conn.Integration.Endpoints))
	for name := range conn.Integration.Endpoints {
		names = append(names, name)
	}
	return names
}
