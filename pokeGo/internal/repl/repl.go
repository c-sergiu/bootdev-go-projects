package repl

import (
	"bufio"
	"fmt"
	"strings"
)

type CliCommand struct {
	Name string
	Description string
	Callback func([]string) error
}

func NewCliCommand(name, description string, callback func([]string) error) *CliCommand{
	return &CliCommand{
		Name: name,
		Description: description,
		Callback: callback,
	}
}

type CliCommands struct {
	Commands map[string]CliCommand
}

func (cmds *CliCommands) AddCommand(name, description string, callback func(options []string) error ) {
	cmd := NewCliCommand(name, description, callback)
	cmds.Commands[name] = *cmd
}

func NewCommands() *CliCommands {
	return &CliCommands{
		Commands: make(map[string]CliCommand),
	}
}

func Loop(scanner *bufio.Scanner, commands CliCommands, cmdDisplPrefix string) {
	for {
		fmt.Print(cmdDisplPrefix)
		if scanner.Scan() {
			text := scanner.Text()
			tokens := cleanInput(text)
			if len(tokens) > 0 {
				cmd, ok := commands.Commands[tokens[0]]
				if ok {
					var t []string
					if len(tokens) > 1 {
						t = tokens[1:]
					} else {
						t = nil
					}
					if err := cmd.Callback(t); err != nil {
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
