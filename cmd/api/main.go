package main

import (
	"fmt"
	"log"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
	"github.com/sanskarchoudhry/pokedex-backend/internal/database"
	"github.com/sanskarchoudhry/pokedex-backend/internal/server"
)

func main() {

	cfg := config.LoadConfig()

	dbService := database.New(cfg.DBUrl)

	defer dbService.Close()

	srv := server.NewServer(cfg, dbService.GetDB())

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}
