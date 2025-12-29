package models

import "time"

type Pokemon struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PokedexID int       `json:"pokedex_id"`
	Name      string    `json:"name"`
	Nickname  string    `json:"nickname,omitempty"`
	Type      string    `json:"type"`
	Height    int       `json:"height"`
	Weight    int       `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}
