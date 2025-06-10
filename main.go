package main

import (
	"log"
	"net/http"

	"fix-ticket-system/config"
	"fix-ticket-system/metrics"
	"fix-ticket-system/models"
	"fix-ticket-system/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ticketService service.TicketServiceInterface

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	config.InitDB()

	// Initialize service
	ticketService = service.NewTicketService()

	// Initialize router
	router := gin.Default()

	// Add Prometheus middleware
	router.Use(metrics.PrometheusMiddleware())

	// Initialize routes
	InitializeRoutes(router)

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// InitializeRoutes initializes the application's routes
func InitializeRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes group
	api := router.Group("/api/v1")
	{
		// Ticket routes
		tickets := api.Group("/tickets")
		{
			tickets.POST("/", createTicket)
			tickets.GET("/", getTickets)
			tickets.GET("/:id", getTicket)
			tickets.PUT("/:id", updateTicket)
			tickets.DELETE("/:id", deleteTicket)
		}
	}
}

// Handler functions
func createTicket(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		CreatedBy   string `json:"created_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		metrics.ErrorTotal.WithLabelValues("invalid_input").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := ticketService.CreateTicket(input.Title, input.Description, input.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func getTickets(c *gin.Context) {
	tickets, err := ticketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func getTicket(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("invalid_id").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := ticketService.GetTicket(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func updateTicket(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("invalid_id").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input struct {
		Title       string          `json:"title" binding:"required"`
		Description string          `json:"description" binding:"required"`
		Status      models.Status   `json:"status" binding:"required"`
		Priority    models.Priority `json:"priority" binding:"required"`
		AssignedTo  string          `json:"assigned_to"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		metrics.ErrorTotal.WithLabelValues("invalid_input").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := ticketService.UpdateTicket(id, input.Title, input.Description, input.Status, input.Priority, input.AssignedTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func deleteTicket(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("invalid_id").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	if err := ticketService.DeleteTicket(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}
