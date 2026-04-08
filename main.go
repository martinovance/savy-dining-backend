package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Reservation struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	TableSize int       `json:"table_size" binding:"required,min=1"`
	Time      time.Time `json:"time" binding:"required"`
	Status    string    `json:"status"`
}

var db *sql.DB

func main() {
	// Initialize Database
	initDB()
	defer db.Close()

	r := gin.Default()

	// Welcome Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Savy Dining API",
			"version": "1.1.0",
			"status":  "running",
			"db":      "connected",
		})
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "error": "database unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "up", "database": "connected"})
	})

	// API Group
	api := r.Group("/api/v1")
	{
		api.POST("/reservations", createReservation)
		api.GET("/reservations", getReservations)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func initDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Printf("Warning: Could not connect to database: %v", err)
	} else {
		log.Println("Successfully connected to database")
		createTables()
	}
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS reservations (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		table_size INT NOT NULL,
		reservation_time TIMESTAMP NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating table: %v", err)
	}
}

func createReservation(c *gin.Context) {
	var newRes Reservation
	if err := c.ShouldBindJSON(&newRes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRes.ID = uuid.New().String()
	newRes.Status = "confirmed"

	query := `INSERT INTO reservations (id, name, email, table_size, reservation_time, status) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, newRes.ID, newRes.Name, newRes.Email, newRes.TableSize, newRes.Time, newRes.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reservation"})
		return
	}

	c.JSON(http.StatusCreated, newRes)
}

func getReservations(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email, table_size, reservation_time, status FROM reservations ORDER BY reservation_time DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}
	defer rows.Close()

	var results []Reservation
	for rows.Next() {
		var res Reservation
		if err := rows.Scan(&res.ID, &res.Name, &res.Email, &res.TableSize, &res.Time, &res.Status); err != nil {
			continue
		}
		results = append(results, res)
	}

	if results == nil {
		results = []Reservation{}
	}

	c.JSON(http.StatusOK, results)
}
