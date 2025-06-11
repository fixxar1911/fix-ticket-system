package routes

import (
	"net/http"

	"fix-ticket-system/middleware"
	"fix-ticket-system/models"
	"fix-ticket-system/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminRoutes struct {
	userService *service.UserService
	auth        *middleware.AuthMiddleware
}

func NewAdminRoutes(userService *service.UserService, auth *middleware.AuthMiddleware) *AdminRoutes {
	return &AdminRoutes{
		userService: userService,
		auth:        auth,
	}
}

func (r *AdminRoutes) Register(router *gin.Engine) {
	admin := router.Group("/api/v1/admin")
	admin.Use(r.auth.RequireAuth(), r.auth.RequireAdmin())

	admin.POST("/users", r.createUser)
	admin.GET("/users", r.listUsers)
	admin.GET("/users/:id", r.getUser)
	admin.PUT("/users/:id", r.updateUser)
	admin.DELETE("/users/:id", r.deleteUser)
}

func (r *AdminRoutes) createUser(c *gin.Context) {
	var input struct {
		Email    string      `json:"email" binding:"required,email"`
		Password string      `json:"password" binding:"required,min=6"`
		Role     models.Role `json:"role" binding:"required,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := r.userService.CreateUser(input.Email, input.Password, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (r *AdminRoutes) listUsers(c *gin.Context) {
	users, err := r.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (r *AdminRoutes) getUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := r.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *AdminRoutes) updateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input struct {
		Email string      `json:"email" binding:"required,email"`
		Role  models.Role `json:"role" binding:"required,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := r.userService.UpdateUser(id, input.Email, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *AdminRoutes) deleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := r.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
