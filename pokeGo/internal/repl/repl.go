package repl

import (
	"bufio"
	"os"
	"fmt"
	"strings"
)

type REPL[T any] struct {
	Commands map[string]Command[T]
	Config T
	Prompt string
}

func (r *REPL[T]) Register(cmd Command[T]) {
	r.Commands[cmd.Name] = cmd
}

func (r *REPL[T]) Run() {
	scanner := bufio.NewScanner(os.Stdout)

	for {
		fmt.Print(r.Prompt)
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		tokens := strings.Fields(input)
		if len(tokens) < 1 {
			break
		}
		cmd, ok := r.Commands[tokens[0]]
		if ok {
			args := tokens[1:]
			if err := cmd.Exec(r.Config, args); err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Printf("command: %s does not exit\n", tokens[0])
			break
		}
	}

}

func NewREPL[T any](config T, prompt string) *REPL[T] {
	r := &REPL[T]{
		Commands: make(map[string]Command[T]),
		Config: config,
		Prompt: prompt,
	}
	r.Register(Command[T]{
		Name: "help",
		Desc: "Displays the commands available",
		Exec: func(cfg T, args []string) error {
			for _, cmd := range r.Commands {
				fmt.Printf("%s: %s\n", cmd.Name, cmd.Desc)
			}
			return nil
		},
	})
	r.Register(Command[T]{
		Name: "exit",
		Desc: "Exit the program",
		Exec: func(cfg T, args []string) error {
			fmt.Println("Goodbye!")
			os.Exit(0)
			return nil
		},
	})
	return r
}


type Command[T any] struct {
	Name string
	Desc string
	Exec func(cfg T, args []string) error
}
