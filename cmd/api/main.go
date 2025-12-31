package main

import (
	"log/slog"
	"os"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
	"github.com/sanskarchoudhry/pokedex-backend/internal/database"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
	"github.com/sanskarchoudhry/pokedex-backend/internal/server"
	"github.com/sanskarchoudhry/pokedex-backend/internal/service"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Pokedex Backend", "env", "dev", "version", "1.0.0")

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
	srv := server.NewServer(cfg, logger, authSvc, pokeService)

	slog.Info("Server running", "port", cfg.Port)
	if err := srv.Start(); err != nil {
		slog.Error("Server crashed", "error", err)
		os.Exit(1)
	}
}
