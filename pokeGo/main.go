package main

import (
	"fmt"
	"bufio"
	"os"
	"io"
	"time"
	"encoding/json"
	"github.com/c-sergiu/pokego/internal/cache"
	"github.com/c-sergiu/pokego/internal/repl"
	"math/rand"
)

func main() {
	var config = NewConfig("https://pokeapi.co/api/v2/location-area/")
	cache := cache.NewCache(60 * time.Second)
	pokedex := NewPokedex()
	commands := repl.NewCommands()

	commands.AddCommand("exit", "Exit the Pokedex", func(options []string) error {
		return CommandExit()
	})
	commands.AddCommand("help", "Display a help message", func(options []string) error {
		return CommandHelp(commands.Commands)
	})
	commands.AddCommand("map", "Display next area locations", func(options []string) error {
		return CommandMap(cache, config)
	})
	commands.AddCommand("mapb", "Displays previous area locations", func(options []string) error {
		return CommandMapb(cache, config)
	})
	commands.AddCommand("explore", "Explore area locations", func(options []string) error {
		if options == nil {
			return fmt.Errorf("Please provide the area's name")
		}
		return CommandExplore(cache, options[0])
	})
	commands.AddCommand("catch", "Catch a pokemon", func(options []string) error {
		if options == nil {
			return fmt.Errorf("Please provide the pokemon's name")
		}
		return CommandCatch(pokedex, cache, options[0])
	})
	commands.AddCommand("inspect", "Inspect caught pokemons", func(options []string) error {
		if options == nil {
			return fmt.Errorf("Please provide the pokemon's name")
		}
		return CommandInspect(pokedex, options[0])
	})
	commands.AddCommand("pokedex", "View caught pokemons", func(options []string) error {
		pokedex.Print()
		return nil
	})

	scanner := bufio.NewScanner(os.Stdin)
	repl.Loop(scanner, *commands, "Pokedex > ")
}
