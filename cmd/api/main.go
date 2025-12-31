package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
	"github.com/sanskarchoudhry/pokedex-backend/internal/database"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
	"github.com/sanskarchoudhry/pokedex-backend/internal/server"
	"github.com/sanskarchoudhry/pokedex-backend/internal/service"
)

func main() {
	// 1. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Database
	dbService := database.New(cfg.DBUrl)

	// 4. Wiring
	userRepo := repository.NewUserRepository(dbService.GetDB())
	tokenRepo := repository.NewTokenRepository(dbService.GetDB())
	pokeRepo := repository.NewPokemonRepository(dbService.GetDB())

	authSvc := service.NewAuthService(userRepo, tokenRepo)
	pokeSvc := service.NewPokemonService(pokeRepo)

	srv := server.NewServer(cfg, logger, authSvc, pokeSvc)

	// 5. Start Server in a Goroutine (Background)
	go func() {
		logger.Info("Server starting", "port", cfg.Port)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// 6. Wait for Shutdown Signal
	// We create a channel that listens for OS signals (Ctrl+C, Docker Stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// This line BLOCKS until a signal is received
	<-quit

	logger.Info("Server shutting down...")

	// 7. Graceful Shutdown
	// Create a context with a 5-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	// 8. Close Database Connections
	// Now that no more requests are coming in, it is safe to close the DB.
	if err := dbService.Close(); err != nil {
		logger.Error("Error closing database", "error", err)
	}

	logger.Info("Server exited properly")
}
