package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sanskarchoudhry/pokedex-backend/internal/models"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
)

// Define custom errors so the Handler knows what status code to send
// (e.g., ErrInvalidInput -> 400 Bad Request)
var (
	ErrInvalidInput = errors.New("invalid input data")
)

type PokemonService interface {
	Create(ctx context.Context, userId int, pokedexId int, name, nickname, pokemonType string, height, weight int) (*models.Pokemon, error)
	List(ctx context.Context, userId int) ([]models.Pokemon, error)
}

type pokemonService struct {
	pokemonRepo repository.PokemonRepository
}

func NewPokemonService(repo repository.PokemonRepository) PokemonService {
	return &pokemonService{
		pokemonRepo: repo,
	}
}

func (p *pokemonService) Create(ctx context.Context, userId int, pokedexId int, name, nickname, pokemonType string, height, weight int) (*models.Pokemon, error) {

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
	}
	if pokedexId <= 0 {
		return nil, fmt.Errorf("%w: pokedex_id must be positive", ErrInvalidInput)
	}
	if height <= 0 || weight <= 0 {
		return nil, fmt.Errorf("%w: height and weight must be positive", ErrInvalidInput)
	}

	// Normalize inputs (e.g., trim whitespace)
	cleanName := strings.TrimSpace(name)
	cleanNickname := strings.TrimSpace(nickname)

	// Rule: If nickname is empty, default to the Pokemon Name
	if cleanNickname == "" {
		cleanNickname = cleanName
	}

	newPokemon := &models.Pokemon{
		UserID:    userId,
		PokedexID: pokedexId,
		Name:      cleanName,
		Nickname:  cleanNickname,
		Type:      pokemonType,
		Height:    height,
		Weight:    weight,
	}

	if err := p.pokemonRepo.CreatePokemon(ctx, newPokemon); err != nil {
		return nil, fmt.Errorf("failed to save pokemon: %w", err)
	}

	return newPokemon, nil
}

func (p *pokemonService) List(ctx context.Context, userId int) ([]models.Pokemon, error) {
	// We return a nil slice instead of an error if the user has 0 pokemon,
	// but we must check if the DB failed.
	list, err := p.pokemonRepo.ListPokemonByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to list pokemon: %w", err)
	}

	// If list is nil (from DB), return an empty slice [] instead of null to frontend.
	if list == nil {
		return []models.Pokemon{}, nil
	}

	return list, nil
}
