package main

import (
	"fmt"

	"github.com/sanskarchoudhry/pokedex-backend/internal/server"
)

func main() {
	server := server.NewServer()

	fmt.Println("Server is running on port 8080...")
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
