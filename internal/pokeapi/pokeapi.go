package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type PokeApiClient struct {
	cache *pokecache.Cache
}

type PokeApiLocationArea struct {
	Name string `json:"name"`
}

type PokeApiPokemonAbility struct {
	Name string `json:"name"`
}

type PokeApiPokemonAbilities struct {
	Ability PokeApiPokemonAbility `json:"ability"`
}

type PokeApiPokemonType struct {
	Name string `json:"name"`
}

type PokeApiPokemonTypes struct {
	Type PokeApiPokemonType `json:"type"`
}

type PokeApiPokemon struct {
	Name string `json:"name"`
}

type PokeApiPokemonEncounter struct {
	Pokemon PokeApiPokemon `json:"pokemon"`
}

type PokeApiLocAreaResponse struct {
	Count    int                   `json:"count"`
	Next     string                `json:"next"`
	Previous string                `json:"previous"`
	Results  []PokeApiLocationArea `json:"results"`
}

type PokeApiLocAreaDetailsResponse struct {
	Name              string                    `json:"name"`
	PokemonEncounters []PokeApiPokemonEncounter `json:"pokemon_encounters"`
}

type PokeApiPokemonDetailsResponse struct {
	Name           string                    `json:"name"`
	BaseExperience int                       `json:"base_experience"`
	Types          []PokeApiPokemonTypes     `json:"types"`
	Abilities      []PokeApiPokemonAbilities `json:"abilities"`
}

func NewClient() *PokeApiClient {
	cache := pokecache.NewCache(5 * time.Minute)

	return &PokeApiClient{cache: cache}
}

func (c *PokeApiClient) LocationArea(url string) (*PokeApiLocAreaResponse, error) {
	locAreas, cached := c.cache.Get(url)

	// If not cached - send HTTP request
	if !cached {
		resp, err := http.Get(url)

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("PokeAPI error: %s", err.Error())
		}

		c.cache.Add(url, body)
		locAreas = body
	}

	var pokeApiResponse PokeApiLocAreaResponse
	err := json.Unmarshal(locAreas, &pokeApiResponse)
	if err != nil {
		return nil, err
	}

	return &pokeApiResponse, nil
}

func (c *PokeApiClient) LocationAreaDetails(url string) (*PokeApiLocAreaDetailsResponse, error) {
	locAreas, cached := c.cache.Get(url)

	// If not cached - send HTTP request
	if !cached {
		resp, err := http.Get(url)

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("PokeAPI error: %s", err.Error())
		}

		c.cache.Add(url, body)
		locAreas = body
	}

	var pokeApiResponse PokeApiLocAreaDetailsResponse
	err := json.Unmarshal(locAreas, &pokeApiResponse)
	if err != nil {
		return nil, err
	}

	return &pokeApiResponse, nil
}

func (c *PokeApiClient) PokemonDetails(url string) (*PokeApiPokemonDetailsResponse, error) {
	locAreas, cached := c.cache.Get(url)

	// If not cached - send HTTP request
	if !cached {
		resp, err := http.Get(url)

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("PokeAPI error: %s", err.Error())
		}

		c.cache.Add(url, body)
		locAreas = body
	}

	var pokeApiResponse PokeApiPokemonDetailsResponse
	err := json.Unmarshal(locAreas, &pokeApiResponse)
	if err != nil {
		return nil, err
	}

	return &pokeApiResponse, nil
}
