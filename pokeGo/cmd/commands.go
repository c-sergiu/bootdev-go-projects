package main

import (
	"fmt"
	"math/rand"
	"encoding/json"
	"github.com/c-sergiu/pokego/internal/pokego"
)

func CommandInspect(ctx *pokego.Context, options []string) error {
	if len(options) < 1 {
		return fmt.Errorf("please provide the pokemon's name")
	}
	name := options[0]
	if p, ok := ctx.Dex[name]; ok {
		fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", p.Name, p.Height, p.Weight)
		fmt.Println("Stats:")
		for _, stat := range p.Stats {
			fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range p.Types {
			fmt.Printf(" - %s\n", t.Type.Name)
		}
	} else {
		return fmt.Errorf("%s has not been caught yet!", name)
	}
	return nil
}

func CommandCatch(ctx *pokego.Context, options []string) error {
	if len(options) < 1 {
		return fmt.Errorf("please provide the pokemon's name")
	}
	name := options[0]
	data, err := ctx.Client.GetPokemon(name)
	if err != nil {return err}

	var p pokego.Pokemon
	if err := json.Unmarshal(data, &p); err != nil {return err}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", p.Name)
	if rand.Intn(700) > p.BaseExperience {
		fmt.Printf("%s was caught!\n", p.Name)
		ctx.Dex[name] = p
	} else {
		fmt.Printf("%s escaped!\n", p.Name)
	}

	return nil
}

func CommandExplore(ctx *pokego.Context, options []string) error {
	if len(options) < 1 {
		return fmt.Errorf("please provide the map name")
	}
	data, err := ctx.Client.GetMap(options[0])
	if err != nil {return err}

	var encounters pokego.Encounters
	if err := json.Unmarshal(data, &encounters); err != nil {return err}

	for _, e := range encounters.PokemonEncounters {
		fmt.Println(e.Pokemon.Name)
	}

	return nil
}

func CommandMap(ctx *pokego.Context, options []string) error {
	var url string
	if ctx.Nav.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	url = *ctx.Nav.Next
	data, err := ctx.Client.GetMaps(url)
	if err != nil {return err}

	var r pokego.LocAreaResult
	if err := json.Unmarshal(data, &r); err != nil {return err}

	r.DisplayResults()
	ctx.Nav.Update(r.Prev, r.Next)
	return nil
}

func CommandMapb(ctx *pokego.Context, options []string) error {
	var url string
	if ctx.Nav.Prev == nil {
		return fmt.Errorf("you're on the first page")
	}
	url = *ctx.Nav.Prev
	data, err := ctx.Client.GetMaps(url)
	if err != nil { return err }

	var r pokego.LocAreaResult
	if err := json.Unmarshal(data, &r); err != nil {return err}

	r.DisplayResults()
	ctx.Nav.Update(r.Prev, r.Next)
	return nil
}

func CommandPokedex(ctx *pokego.Context, options []string) error {
	for _, p := range ctx.Dex {
		fmt.Printf(" - %s\n", p.Name)
	}
	return nil
}
