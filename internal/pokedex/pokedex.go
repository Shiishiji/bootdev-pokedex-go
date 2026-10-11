package pokedex

import (
	"fmt"
)

type CaughtPokemon struct {
	Name string
}

type Pokedex struct {
	caught map[string]CaughtPokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{map[string]CaughtPokemon{}}
}

func (p *Pokedex) Register(pokemonName string) error {
	_, alreadyCaught := p.caught[pokemonName]
	if alreadyCaught {
		return fmt.Errorf("Pokemon %s already in pokedex", pokemonName)
	}

	p.caught[pokemonName] = CaughtPokemon{Name: pokemonName}
	return nil
}

func (p *Pokedex) Has(pokemonName string) bool {
	_, isCaught := p.caught[pokemonName]
	return isCaught
}

func (p *Pokedex) List() []CaughtPokemon {
	pokemons := []CaughtPokemon{}
	for _, p := range p.caught {
		pokemons = append(pokemons, p)
	}

	return pokemons
}
