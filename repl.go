package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/KrisQ/verbose-bassoon/internal/pokeapi"
)

type Config struct {
	pokeapiClient pokeapi.Client
	Next          *string
	Previous      *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
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
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(config *Config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func commandMap(config *Config) error {
	locationAreas, err := config.pokeapiClient.GetLocations(config.Next)
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

func commandMapb(config *Config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	locationAreas, err := config.pokeapiClient.GetLocations(config.Previous)
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

func startRepl(config *Config) {
	scanner := bufio.NewScanner(os.Stdin)
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
		err := command.callback(config)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
