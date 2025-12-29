package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sanskarchoudhry/pokedex-backend/internal/models"
)

type PokemonRepository interface {
	CreatePokemon(ctx context.Context, p *models.Pokemon) error
	ListPokemonByUserID(ctx context.Context, userID int) ([]models.Pokemon, error)
}

type postgresPokemonRepository struct {
	db *sql.DB
}

func NewPokemonRepository(db *sql.DB) PokemonRepository {
	return &postgresPokemonRepository{db: db}
}

func (r *postgresPokemonRepository) CreatePokemon(ctx context.Context, p *models.Pokemon) error {
	query := `
		INSERT INTO pokemons (user_id, pokedex_id, name, nickname, type, height, weight)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		p.UserID, p.PokedexID, p.Name, p.Nickname, p.Type, p.Height, p.Weight,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *postgresPokemonRepository) ListPokemonByUserID(ctx context.Context, userID int) ([]models.Pokemon, error) {
	query := `SELECT id, user_id, pokedex_id, name, nickname, type, height, weight, created_at FROM pokemons WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var pokemons []models.Pokemon
	for rows.Next() {
		var p models.Pokemon
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.PokedexID, &p.Name, &p.Nickname,
			&p.Type, &p.Height, &p.Weight, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		pokemons = append(pokemons, p)
	}
	return pokemons, nil
}
