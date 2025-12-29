package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreatePokemonRequest struct {
	PokedexID int    `json:"pokedex_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Nickname  string `json:"nickname"` // Optional
	Type      string `json:"type" binding:"required"`
	Height    int    `json:"height" binding:"required"`
	Weight    int    `json:"weight" binding:"required"`
}

func (s *Server) createPokemonHandler(c *gin.Context) {
	// 1. Get User ID from Context (Set by Auth Middleware)
	// We cast it to int because c.Get returns interface{}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 2. Parse JSON
	var req CreatePokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Call Service
	pokemon, err := s.pokemonService.Create(
		c.Request.Context(),
		userID.(int), // Type assertion: "I promise this is an int"
		req.PokedexID,
		req.Name,
		req.Nickname,
		req.Type,
		req.Height,
		req.Weight,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pokemon"})
		return
	}

	c.JSON(http.StatusCreated, pokemon)
}

func (s *Server) listPokemonHandler(c *gin.Context) {
	// 1. Get User ID from Context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 2. Call Service
	list, err := s.pokemonService.List(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pokemon"})
		return
	}

	// 3. Return JSON
	c.JSON(http.StatusOK, gin.H{"data": list})
}
