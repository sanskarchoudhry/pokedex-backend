package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Simple Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I am alive!"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			// We will write these handlers next!
			auth.POST("/register", s.registerHandler)
			auth.POST("/login", s.loginHandler)
		}
	}

	return r
}

// Stubs to make the compiler happy until we implement logic
func (s *Server) registerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "register endpoint"})
}

func (s *Server) loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "login endpoint"})
}
