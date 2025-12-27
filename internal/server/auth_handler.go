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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) loginHandler(c *gin.Context) {
	var req LoginRequest

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 2. Call Service
	accessToken, refreshToken, err := s.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Security: Don't tell the user exactly what went wrong (user not found vs wrong password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 3. Set Refresh Token in HttpOnly Cookie
	// SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*3600,
		"/",
		"",
		false,
		true,
	)

	// 4. Return Access Token in JSON
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"message":      "Login successful",
	})
}

func (s *Server) refreshHandler(c *gin.Context) {
	// 1. Get the HttpOnly Cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	// 2. Call Service
	newAccessToken, err := s.authService.Refresh(c.Request.Context(), cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// 3. Return the Fresh Access Token
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"message":      "Token refreshed successfully",
	})
}
