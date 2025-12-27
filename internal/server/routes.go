package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.registerHandler)
			auth.POST("/login", s.loginHandler)
			auth.POST("/refresh", s.refreshHandler)
		}

		// Protected Routes
		// We create a new group and apply the Middleware
		protected := v1.Group("/pokedex")
		protected.Use(s.AuthMiddleware())
		{
			protected.GET("/me", func(c *gin.Context) {
				// Retrieve the UserID we set in the middleware
				userID, _ := c.Get("userID")
				c.JSON(http.StatusOK, gin.H{
					"message": "You have accessed a protected route!",
					"user_id": userID,
				})
			})
		}
	}

	return r
}
