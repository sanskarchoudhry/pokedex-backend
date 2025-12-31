package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanskarchoudhry/pokedex-backend/internal/service"
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
	// 1. Context Logging Fields (Traceability)
	// We want every log line here to have "handler=createPokemon"
	log := s.logger.With("handler", "createPokemon")

	userID, exists := c.Get("userID")
	if !exists {
		log.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreatePokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid JSON body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pokemon, err := s.pokemonService.Create(
		c.Request.Context(),
		userID.(int),
		req.PokedexID,
		req.Name,
		req.Nickname,
		req.Type,
		req.Height,
		req.Weight,
	)

	if err != nil {
		// Check if it's a Validation Error or a System Error
		if errors.Is(err, service.ErrInvalidInput) {
			log.Info("Validation failed", "user_id", userID, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Error("Failed to create pokemon", "user_id", userID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	log.Info("Pokemon created", "user_id", userID, "pokemon_id", pokemon.ID, "name", pokemon.Name)
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
