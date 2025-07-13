package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"vietnam-admin-api/handlers"
	"vietnam-admin-api/middleware"
	"vietnam-admin-api/services"
)

const (
	Version         = "1.0.0"
	DefaultPort     = "8100"
	DefaultDataPath = "./data"
)

func main() {
	log.Println("🚀 Starting Vietnam Administrative API Server...")

	// Get configuration from environment
	port := getEnv("PORT", DefaultPort)
	dataPath := getEnv("DATA_PATH", DefaultDataPath)
	ginMode := getEnv("GIN_MODE", "release")

	// Set Gin mode
	gin.SetMode(ginMode)

	// Initialize data service
	dataService := services.NewDataService(dataPath)

	// Load data on startup
	log.Println("📊 Loading administrative data...")
	if err := dataService.LoadData(); err != nil {
		log.Fatalf("❌ Failed to load data: %v", err)
	}

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(dataService, Version)

	// Setup Gin router
	router := setupRouter(apiHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🌐 Server starting on port %s", port)
		log.Printf("📚 Health check: http://localhost:%s/api/v1/health", port)
		log.Printf("📖 API docs: http://localhost:%s/api/v1/stats", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🔄 Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}

func setupRouter(apiHandler *handlers.APIHandler) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Rate limiting (optional - uncomment if needed)
	// router.Use(middleware.RateLimit())

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Province endpoints
		provinces := v1.Group("/provinces")
		{
			provinces.GET("", apiHandler.GetProvinces)
			provinces.GET("/types", apiHandler.GetProvinceTypes)
			provinces.GET("/:code", apiHandler.GetProvince)
			provinces.GET("/:code/wards", apiHandler.GetProvinceWards)
		}

		// Ward endpoints
		wards := v1.Group("/wards")
		{
			wards.GET("", apiHandler.GetWards)
			wards.GET("/types", apiHandler.GetWardTypes)
			wards.GET("/:code", apiHandler.GetWard)
		}

		// Search endpoints
		v1.GET("/search", apiHandler.GlobalSearch)

		// Utility endpoints
		v1.POST("/address/validate", apiHandler.ValidateAddress)
		v1.GET("/health", apiHandler.Health)
		v1.GET("/stats", apiHandler.Stats)

		// Admin endpoints (can be protected with auth middleware)
		admin := v1.Group("/admin")
		// admin.Use(middleware.AdminAuth()) // Uncomment for auth
		{
			admin.POST("/reload", apiHandler.ReloadData)
		}
	}

	// Root endpoints
	router.GET("/health", apiHandler.Health)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Vietnam Administrative API",
			"version": Version,
			"status":  "running",
			"time":    time.Now(),
			"endpoints": gin.H{
				"health":    "/api/v1/health",
				"stats":     "/api/v1/stats",
				"provinces": "/api/v1/provinces",
				"wards":     "/api/v1/wards",
				"search":    "/api/v1/search",
				"validate":  "/api/v1/address/validate",
			},
		})
	})

	// Handle 404
	router.NoRoute(apiHandler.NotFound)

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
