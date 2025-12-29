package main

import (
	"fmt"
	"log"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
	"github.com/sanskarchoudhry/pokedex-backend/internal/database"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
	"github.com/sanskarchoudhry/pokedex-backend/internal/server"
	"github.com/sanskarchoudhry/pokedex-backend/internal/service"
)

func main() {
	// Config
	cfg := config.LoadConfig()

	// Database
	dbService := database.New(cfg.DBUrl)
	defer dbService.Close()

	// Repository

	userRepo := repository.NewUserRepository(dbService.GetDB())
	tokenRepo := repository.NewTokenRepository(dbService.GetDB())
	pokeRepo := repository.NewPokemonRepository(dbService.GetDB())

	// Service
	authSvc := service.NewAuthService(userRepo, tokenRepo)
	pokeService := service.NewPokemonService(pokeRepo)

	// Server
	// We inject the Service into the Server
	srv := server.NewServer(cfg, authSvc, pokeService)

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}
