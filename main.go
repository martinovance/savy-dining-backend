package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Reservation struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	TableSize int       `json:"table_size" binding:"required,min=1"`
	Time      time.Time `json:"time" binding:"required"`
	Status    string    `json:"status"`
}

var reservations = []Reservation{}

func main() {
	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	// API Group
	api := r.Group("/api/v1")
	{
		api.POST("/reservations", createReservation)
		api.GET("/reservations", getReservations)
	}

	r.Run(":8080")
}

func createReservation(c *gin.Context) {
	var newRes Reservation
	if err := c.ShouldBindJSON(&newRes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRes.ID = uuid.New().String()
	newRes.Status = "confirmed"
	reservations = append(reservations, newRes)

	c.JSON(http.StatusCreated, newRes)
}

func getReservations(c *gin.Context) {
	c.JSON(http.StatusOK, reservations)
}
