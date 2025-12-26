package server

import (
	"net/http"
	"time"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
	"github.com/sanskarchoudhry/pokedex-backend/internal/service"
)

type Server struct {
	config      *config.Config
	authService service.AuthService
}

func NewServer(cfg *config.Config, authService service.AuthService) *Server {
	return &Server{
		config:      cfg,
		authService: authService,
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
