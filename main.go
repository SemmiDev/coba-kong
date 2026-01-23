package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables dari file .env jika ada
	// Ini opsional karena di production biasanya env vars di-set langsung
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load konfigurasi aplikasi
	cfg := LoadConfig()

	// Setup logger menggunakan logrus untuk logging yang lebih baik
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Set mode Gin berdasarkan environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inisialisasi router Gin
	router := gin.New()

	// Middleware global - ini akan dijalankan untuk setiap request
	router.Use(Logger(logger)) // Custom logger middleware
	router.Use(gin.Recovery()) // Recovery dari panic
	router.Use(CORS())         // CORS headers

	// Health check endpoint - penting untuk monitoring dan load balancer
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": cfg.ServiceName,
		})
	})

	// Initialize handlers dengan dependencies yang dibutuhkan
	userHandler := NewUserHandler(logger)

	// API routes - grouping untuk versioning
	v1 := router.Group("/api/v1")
	{
		// User endpoints
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Start server dengan graceful shutdown
	logger.Infof("Starting %s on port %s", cfg.ServiceName, cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
