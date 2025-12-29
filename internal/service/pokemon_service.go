package service

import (
	"context"

	"github.com/sanskarchoudhry/pokedex-backend/internal/models"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
)

type PokemonService interface {
	Create(ctx context.Context, userId int, pokedexId int, name string, nickname string, pokemonType string, height int, weight int) (*models.Pokemon, error)
	List(ctx context.Context, userId int) ([]models.Pokemon, error)
}

type pokemonService struct {
	pokemonRepo repository.PokemonRepository
}

func NewPokemonService(pokmonRepo repository.PokemonRepository) PokemonService {
	return &pokemonService{
		pokemonRepo: pokmonRepo,
	}
}

func (p *pokemonService) Create(ctx context.Context, userId int, pokedexId int, name string, nickname string, pokemonType string, height int, weight int) (*models.Pokemon, error) {
	newPokemon := &models.Pokemon{UserID: userId, PokedexID: pokedexId, Name: name, Nickname: nickname, Height: height, Weight: weight, Type: pokemonType}

	if err := p.pokemonRepo.CreatePokemon(ctx, newPokemon); err != nil {
		return nil, err
	}

	responseMon := *newPokemon
	return &responseMon, nil

}

func (p *pokemonService) List(ctx context.Context, userId int) ([]models.Pokemon, error) {
	if pokemonList, err := p.pokemonRepo.ListPokemonByUserID(ctx, userId); err != nil {
		return nil, err
	} else {
		return pokemonList, nil
	}
}
