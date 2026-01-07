package main

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"github.com/c-sergiu/pokego/internal/pokego"
)

type CliCommand struct {
	Name string
	Description string
	Callback func(*pokego.Context, []string) error
}

func main() {
	ctx := pokego.NewContext()
	cmd := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("pokego > ")
		if scanner.Scan() {
			text := scanner.Text()
			tokens := cleanInput(text)
			if len(tokens) > 0 {
				cmd, ok := cmd[tokens[0]]
				if ok {
					if err := cmd.Callback(ctx, tokens); err != nil {
						fmt.Printf("%v\n", err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}

func cleanInput(text string) []string {
	out := []string{}
	var curr string
	for i := range text {
		if text[i] == ' ' {
			if len(curr) > 0 {
				out = append(out, strings.ToLower(curr))
				curr = ""
			}
		}else {
			curr += string(text[i])
		}
		if i == len(text) -1 {
			if len(curr) > 0 {
				out = append(out, strings.ToLower(curr))
				curr = ""
			}
		}
	}
	return out
}

func getCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": CliCommand{
			Name: "exit",
			Description: "Exit the Pokedex",
			Callback: CommandExit,
		},
		"help": CliCommand{
			Name: "help", 
			Description: "Display a help message", 
			Callback: CommandHelp,
		},
		"map": CliCommand{
			Name: "map",
			Description: "Display next area locations", 
			Callback: CommandMap,
		},
		"mapb": CliCommand{
			Name:"mapb", 
			Description:"Displays previous area locations", 
			Callback: CommandMapb,
		},
		"explore": CliCommand{
			Name: "explore", 
			Description: "Explore area locations", 
			Callback: CommandExplore,
		},
		"catch": CliCommand{
			Name: "catch", 
			Description: "Catch a pokemon", 
			Callback: CommandCatch,
		},
		"inspect": CliCommand{
			Name: "inspect",
			Description: "Inspect caught pokemons", 
			Callback: CommandInspect, 
		},
		"pokedex": CliCommand{
			Name: "pokedex", 
			Description: "View caught pokemons",
			Callback: CommandPokedex,
		},
	}
}
