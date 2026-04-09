package main

import (
	"log"
	"os"
	"martinovance/savy-dining-backend/internal/repository"
	"martinovance/savy-dining-backend/internal/api"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	repo := repository.NewRepository(db)
	if err := repo.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := gin.Default()
	
	// API V1 Routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})
		v1.GET("/menu", api.GetMenuHandler(repo))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
