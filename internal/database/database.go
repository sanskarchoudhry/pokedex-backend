package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Service represents the database service
type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *sql.DB
}

type service struct {
	db *sql.DB
}

var (
	dbInstance *service
)

// New initializes the database connection
func New(connectionString string) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Process terminated: Cannot connect to DB: %v", err)
	}

	dbInstance = &service{
		db: db,
	}

	fmt.Println("Connected to PostgreSQL Database!")
	return dbInstance
}

// GetDB returns the underlying sql.DB instance
func (s *service) GetDB() *sql.DB {
	return s.db
}

// Close closes the database connection
func (s *service) Close() error {
	log.Printf("Disconnected from Database")
	return s.db.Close()
}

// Health checks the health of the database connection
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log fatal error
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	stats["open_connections"] = fmt.Sprintf("%d", s.db.Stats().OpenConnections)
	return stats
}
