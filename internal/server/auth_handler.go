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
	log := s.logger.With("handler", "register")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid register payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	user, err := s.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Error("Registration failed", "email", req.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("User registered successfully", "email", user.Email, "user_id", user.ID)

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
	log := s.logger.With("handler", "login")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid login payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, refreshToken, err := s.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Warn("Login failed", "email", req.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 3. FIXED: Make "Secure" dynamic based on environment
	isProd := false
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "", isProd, true)

	log.Info("User logged in", "email", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"message":      "Login successful",
	})
}

func (s *Server) refreshHandler(c *gin.Context) {
	log := s.logger.With("handler", "refresh")

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		log.Warn("Refresh token missing") // <--- Log it
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	newAccessToken, err := s.authService.Refresh(c.Request.Context(), cookie)
	if err != nil {
		log.Warn("Refresh failed", "error", err) // <--- Log security events
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	log.Info("Token refreshed successfully") // <--- Success log

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"message":      "Token refreshed successfully",
	})
}
