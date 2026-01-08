package main

import (
	"github.com/c-sergiu/pokego/internal/pokego"
	"github.com/c-sergiu/pokego/internal/repl"
)

func main() {
	ctx := pokego.NewContext()
	repl := repl.NewREPL[*pokego.Context](ctx, "pokego > ")
	registerCommands(repl)
	repl.Run()
}

func registerCommands(r *repl.REPL[*pokego.Context]) {
	r.Register(repl.Command[*pokego.Context]{
		Name: "map",
		Desc: "Display next area locations", 
		Exec: CommandMap,
	})
	r.Register(repl.Command[*pokego.Context]{
		Name:"mapb", 
		Desc:"Displays previous area locations", 
		Exec: CommandMapb,
	})
	r.Register(repl.Command[*pokego.Context]{
		Name: "explore", 
		Desc: "Explore area locations", 
		Exec: CommandExplore,
	})
	r.Register(repl.Command[*pokego.Context]{
		Name: "catch", 
		Desc: "Catch a pokemon", 
		Exec: CommandCatch,
	})
	r.Register(repl.Command[*pokego.Context]{
		Name: "inspect",
		Desc: "Inspect caught pokemons", 
		Exec: CommandInspect, 
	})
	r.Register(repl.Command[*pokego.Context]{
		Name: "pokedex", 
		Desc: "View caught pokemons",
		Exec: CommandPokedex,
	})
}
