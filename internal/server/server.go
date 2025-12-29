package server

import (
	"net/http"
	"time"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
	"github.com/sanskarchoudhry/pokedex-backend/internal/service"
)

type Server struct {
	config         *config.Config
	authService    service.AuthService
	pokemonService service.PokemonService
}

func NewServer(cfg *config.Config, authService service.AuthService, pokeSvc service.PokemonService) *Server {
	return &Server{
		config:         cfg,
		authService:    authService,
		pokemonService: pokeSvc,
	}
}

func (s *Server) Start() error {
	httpServer := &http.Server{
		Addr:         s.config.Port,
		Handler:      s.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return httpServer.ListenAndServe()
}
