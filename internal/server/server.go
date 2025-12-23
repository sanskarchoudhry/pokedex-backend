package server

import "net/http"

type Server struct {
	port string
	// db *sql.DB  <-- We will add this tomorrow!
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: ":8080",
	}

	// Declare Server config
	server := &http.Server{
		Addr:    NewServer.port,
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
