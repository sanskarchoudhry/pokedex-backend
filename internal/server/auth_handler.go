package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanskarchoudhry/pokedex-backend/internal/utils"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (s *Server) registerHandler(c *gin.Context) {
	var req RegisterRequest

	// 1. Validation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 2. Hash the password
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// 3. (TODO) Database Insert
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"debug_info": gin.H{
			"email":       req.Email,
			"hashed_pass": hashedPwd,
		},
	})
}
