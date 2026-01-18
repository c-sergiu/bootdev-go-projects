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
		if err := r.HandleCommand(tokens); err != nil {
			fmt.Println(err)
		}
	}

}
func (r *REPL[T]) HandleCommand(args []string) error {
	args_l := len(args)
	if args_l < 1 {
		return fmt.Errorf("missing command")
	}
	cmd, ok := r.Commands[args[0]]
	if ok {
		if args_l - 1 < cmd.MinArgs {
			return fmt.Errorf("Not enaugh args")
		}
		if err := cmd.Exec(r.Config, args[1:]); err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		return fmt.Errorf("command: %s does not exit\n", args[0])
	}
	return nil
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
		MinArgs: 0,
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
		MinArgs: 0,
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
	MinArgs int
	Exec func(cfg T, args []string) error
}
