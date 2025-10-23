package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/KrisQ/verbose-bassoon/internal/pokeapi"
	"github.com/KrisQ/verbose-bassoon/internal/pokecache"
)

type Config struct {
	pokeapiClient pokeapi.Client
	Next          *string
	Previous      *string
	cache         pokecache.Cache
	pokedex       dex
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

type dex struct {
	caught map[string]pokeapi.Pokemon
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
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
			description: "Display next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore area and find pokemons",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon - difficulty based on pokemon's experience",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "look up pokemon in pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "list caught pokemons",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(config *Config, _ ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, _ ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func commandMap(config *Config, _ ...string) error {
	locationAreas, err := config.pokeapiClient.GetLocations(&config.cache, config.Next)
	if err != nil {
		return err
	}
	for _, loc := range locationAreas.Results {
		fmt.Println(loc.Name)
	}
	config.Previous = locationAreas.Previous
	config.Next = locationAreas.Next
	return nil
}

func commandMapb(config *Config, _ ...string) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	locationAreas, err := config.pokeapiClient.GetLocations(&config.cache, config.Previous)
	if err != nil {
		return err
	}
	for _, loc := range locationAreas.Results {
		fmt.Println(loc.Name)
	}
	config.Previous = locationAreas.Previous
	config.Next = locationAreas.Next
	return nil
}

func commandExplore(config *Config, args ...string) error {
	if args[0] == "" {
		fmt.Println("you need to provide a location")
		return nil
	}
	locationPokemons, err := config.pokeapiClient.GetLocationPokemons(&config.cache, args[0])
	if err != nil {
		return err
	}
	for _, encounter := range locationPokemons.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config, args ...string) error {
	if args[0] == "" {
		fmt.Println("you need to provide a pokemon name")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	pokemon, err := config.pokeapiClient.GetPokemon(&config.cache, args[0])
	if err != nil {
		return err
	}
	base := 1000
	pokemonExp := pokemon.BaseExperience
	random := rand.Intn(1000)
	if base-pokemonExp-random > 500 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		config.pokedex.caught[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(config *Config, args ...string) error {
	if args[0] == "" {
		fmt.Println("you need to provide a location")
		return nil
	}
	pokemon, ok := config.pokedex.caught[args[0]]
	if !ok {
		fmt.Println("you can only inspect pokemon you've caught")
		return nil
	}
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("\t-%v: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("\t-%v\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(config *Config, args ...string) error {
	fmt.Println("Your pokedex:")
	for k := range config.pokedex.caught {
		fmt.Printf("\t-%v\n", k)
	}
	return nil
}

func startRepl(config *Config) {
	scanner := bufio.NewScanner(os.Stdin)
	config.cache = *pokecache.NewCache(5 * time.Second)
	config.pokedex.caught = make(map[string]pokeapi.Pokemon) // Initialize here
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}
		command, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		err := command.callback(config, args...)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
