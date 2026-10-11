package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokedex"
	"strings"
)

type config struct {
	pokeapi          *pokeapi.PokeApiClient
	pokedex          *pokedex.Pokedex
	nextUrl          string
	previousUrl      string
	input            []string
	commandsRegistry map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func main() {
	config := &config{
		pokeapi:     pokeapi.NewClient(),
		pokedex:     pokedex.NewPokedex(),
		nextUrl:     "",
		previousUrl: "",
		commandsRegistry: map[string]cliCommand{
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"map": {
				name:        "map",
				description: "Get next 20 location areas",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Get previous 20 location areas",
				callback:    commandMapb,
			},
			"explore": {
				name:        "explore",
				description: "Get list of pokemon in the area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Throw pokeball at a pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect pokemon information",
				callback:    commandInspect,
			},
			"pokedex": {
				name:        "pokedex",
				description: "List catched pokemons",
				callback:    commandPokedex,
			},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		config.input = cleanInput(scanner.Text())
		if len(config.input) < 1 {
			fmt.Println("No command provided")
			continue
		}

		command := config.input[0]

		commandReference, commandExists := config.commandsRegistry[command]
		if !commandExists {
			fmt.Println("Unknown command")
			continue
		}

		err := commandReference.callback(config)
		if err != nil {
			fmt.Errorf("Command failed with message: %s", err.Error())
		}
	}
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range config.commandsRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(config *config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if config.nextUrl != "" {
		url = config.nextUrl
	}

	response, err := config.pokeapi.LocationArea(url)
	if err != nil {
		return err
	}

	config.nextUrl = response.Next
	config.previousUrl = response.Previous

	for _, locationArea := range response.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}

	return nil
}

func commandMapb(config *config) error {
	if config.previousUrl == "" {
		fmt.Print("you're on the first page\n")
		return nil
	}

	url := config.previousUrl

	response, err := config.pokeapi.LocationArea(url)
	if err != nil {
		return err
	}

	config.nextUrl = response.Next
	config.previousUrl = response.Previous

	for _, locationArea := range response.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}

	return nil
}

func commandExplore(config *config) error {
	if len(config.input) < 2 {
		return fmt.Errorf("Invalid number of parameters")
	}

	locAreaName := config.input[1]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locAreaName)

	response, err := config.pokeapi.LocationAreaDetails(url)
	if err != nil {
		return err
	}

	for _, pokemonEncounter := range response.PokemonEncounters {
		fmt.Printf("%s\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config) error {
	if len(config.input) < 2 {
		return fmt.Errorf("Invalid number of parameters")
	}

	pokemon := config.input[1]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	response, err := config.pokeapi.PokemonDetails(url)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	// Magikarp has 40 BaseExperience and he should almast always be caught
	// So 30.0 / 40.0 gives 75%
	chance := 30.0 / float64(response.BaseExperience) * 100
	chanceBig := big.NewInt(int64(chance))

	fmt.Printf("Catch chance is %d%%\n", chanceBig)

	roll, err := rand.Int(rand.Reader, big.NewInt(100))
	if chanceBig.Cmp(roll) >= 0 {
		fmt.Printf("%s was caught!\n", pokemon)
		err := config.pokedex.Register(pokemon)
		if err != nil {
			fmt.Println(err.Error())
		}

		return nil
	}

	fmt.Printf("%s escaped!\n", pokemon)
	return nil
}

func commandInspect(config *config) error {
	if len(config.input) < 2 {
		return fmt.Errorf("Invalid number of parameters")
	}

	pokemon := config.input[1]

	if !config.pokedex.Has(pokemon) {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	response, err := config.pokeapi.PokemonDetails(url)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", response.Name)
	fmt.Printf("Experience: %d\n", response.BaseExperience)

	fmt.Println("Types:")
	for _, t := range response.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	fmt.Println("Abilities:")
	for _, a := range response.Abilities {
		fmt.Printf(" - %s\n", a.Ability.Name)
	}

	return nil
}

func commandPokedex(config *config) error {
	if len(config.input) < 1 {
		return fmt.Errorf("Invalid number of parameters")
	}

	catchedPokemons := config.pokedex.List()

	fmt.Println("Your Pokedex:")
	for _, p := range catchedPokemons {
		fmt.Printf(" - %s\n", p.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	parsedText := strings.ToLower(strings.TrimSpace(text))
	parsedTextLength := len(parsedText)

	result := []string{}
	buffer := ""
	for i, c := range parsedText {
		if c == ' ' {
			if buffer != "" {
				result = append(result, buffer)
			}

			buffer = ""
			continue
		}

		buffer = buffer + string(c)
		if i+1 == parsedTextLength {
			result = append(result, buffer)
			buffer = ""
		}
	}

	return result
}
