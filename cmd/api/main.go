package main

import (
	"log"
	"os"
	"github.com/martinovance/savy-dining-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// For debugging, log that we are starting
	log.Println("Initializing Savy Dining Backend...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("CRITICAL: DATABASE_URL environment variable is missing")
		// Don't exit immediately so the log can be captured
		os.Exit(1)
	}

	log.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("CRITICAL: Failed to connect to database: %v", err)
		os.Exit(1)
	}

	repo := repository.NewRepository(db)
	log.Println("Running auto-migrations...")
	if err := repo.AutoMigrate(); err != nil {
		log.Printf("CRITICAL: Failed to run migrations: %v", err)
		os.Exit(1)
	}

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	r := gin.Default()
	
	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Root welcome route (added back for basic connectivity check)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Savy Dining API",
			"status":  "running",
		})
	})

	// API V1 Routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "healthy",
				"version": "1.0.2",
				"db":      "connected",
			})
		})
		v1.GET("/menu", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Menu endpoint active"})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Bind to 0.0.0.0 to allow external traffic
	log.Printf("Starting server on 0.0.0.0:%s", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Printf("CRITICAL: Failed to start server: %v", err)
		os.Exit(1)
	}
}
