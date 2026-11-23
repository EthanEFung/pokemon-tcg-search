package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethanefung/pokemon-tcg-search/internal"
	"github.com/jackc/pgx/v5"
)

type Entry struct {
	EntryNumber    int `json:"entry_number"`
	PokemonSpecies struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon_species"`
}

type Body struct {
	PokemonEntries []Entry `json:"pokemon_entries"`
}

func main() {
	endpoint := "http://pokeapi.co/api/v2/pokedex/1" // national dex
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatalf("request err: %v\n", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("call err: %v\n", err)

	}
	defer res.Body.Close()

	var body Body
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Fatalf("call err: %v\n", err)
	}

	for _, entry := range body.PokemonEntries {
		fmt.Println(entry.PokemonSpecies.Name)
	}

	conn, err := pgx.Connect(ctx, "user=ethanfung dbname=pokemon-tcg-search sslmode=verify-full")
	if err != nil {
		log.Fatalf("connection err: %v\n", err)
	}
	defer conn.Close(ctx)
	queries := internal.New(conn)

	params := []internal.CreatePokemonParams{}

	for _, entry := range body.PokemonEntries {
		params = append(params, internal.CreatePokemonParams{
			ID:   int32(entry.EntryNumber),
			Name: entry.PokemonSpecies.Name,
		})
	}

	byt, err := queries.CreatePokemon(ctx, params)
	if err != nil {
		log.Fatalf("insert err: %v\n", err)
	}
	fmt.Printf("%d rows inserted", byt)
}
