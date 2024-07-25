package main

import (
	"fmt"
	pokedex "gokedex/internal/http_get"
	"log"
)

func main() {
	fmt.Println("Type the name of a pokemon")
	pokemonName := ""
	fmt.Scanln(&pokemonName)
	res, err := pokedex.GetPokemon(pokemonName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
