package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (s *Server) registerHandler(c *gin.Context) {
	var req RegisterRequest

	// 1. Validation (HTTP Layer)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 2. Call Service (Business Logic Layer)
	// We pass c.Request.Context() so if the user disconnects, the DB query stops.
	user, err := s.authService.Register(c.Request.Context(), req.Email, req.Password)

	// 3. Handle Errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
