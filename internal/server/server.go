package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/sanskarchoudhry/pokedex-backend/internal/config"
)

type Server struct {
	config *config.Config
	db     *sql.DB
}

// NewServer receives dependencies (DI) instead of creating them
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		config: cfg,
		db:     db,
	}
}

// Start handles the HTTP server lifecycle
func (s *Server) Start() error {
	httpServer := &http.Server{
		Addr:         s.config.Port,
		Handler:      s.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return httpServer.ListenAndServe()
}
