package main

import (
	"fmt"
	"os"
	"github.com/c-sergiu/pokego/internal/repl"
	"github.com/c-sergiu/pokego/internal/cache"
	"math/rand"
	"encoding/json"
	"net/http"
	"io"
)

func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(commands map[string]repl.CliCommand) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.Name, c.Description)
	}
	return nil
}

func CommandInspect(dex *Pokedex, name string) error {
	if p, ok := dex.Dex[name]; ok {
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

func CommandCatch(dex *Pokedex, cache *cache.Cache, name string) error {
	p, err := getPokemon(cache, name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", p.Name)
	if rand.Intn(700) > p.BaseExperience {
		fmt.Printf("%s was caught!\n", p.Name)
		dex.Add(p)
	} else {
		fmt.Printf("%s escaped!\n", p.Name)
	}
	return nil
}


func CommandExplore(cache *cache.Cache, mapName string) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	url := baseUrl + mapName
	var bytes []byte
	p, ok := cache.Get(url)
	if ok {
		bytes = p
	} else {
		if err := getResBody(url, &bytes); err != nil { return err }
		cache.Add(url, bytes)
	}
	var encounters Encounters
	if err := json.Unmarshal(bytes, &encounters); err != nil {
		return fmt.Errorf("Error while unmarshaling encounters: %v", err)
	}
	for _, e := range encounters.PokemonEncounters {
		fmt.Println(e.Pokemon.Name)
	}
	return nil
}

func CommandMap(cache *cache.Cache, conf *NavConfig) error {
	var url string
	if conf.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	url = *conf.Next
	r, err := getMap(url, cache)
	if err != nil {return err}
	r.DisplayResults()
	conf.Update(r.Prev, r.Next)
	return nil
}

func CommandMapb(cache *cache.Cache, conf *NavConfig) error {
	var url string
	if conf.Prev == nil {
		return fmt.Errorf("you're on the first page")
	}
	url = *conf.Prev
	r, err := getMap(url, cache)
	if err != nil { return err }
	r.DisplayResults()
	conf.Update(r.Prev, r.Next)
	return nil
}

func getPokemon(cache *cache.Cache, name string) (Pokemon, error) {
	empty := Pokemon{}
	if len(name) == 0 || name == "" {
		return empty, fmt.Errorf("invalid pokemon name")
	}
	base := "https://pokeapi.co/api/v2/pokemon/"
	url := base + name
	var bytes []byte

	if e, ok := cache.Get(url); ok {
		bytes = e
	}else {
		if err := getResBody(url, &bytes); err != nil { return empty, err }
		cache.Add(url, bytes)
	}
	var pokemon Pokemon
	if err := json.Unmarshal(bytes, &pokemon); err != nil {
		return empty, fmt.Errorf("Error while unmarshalling pokemon: %v", err)
	}
	return pokemon, nil
}

func getMap(url string, cache *cache.Cache) (*LocAreaResult, error) {
	var bytes []byte
	e, ok := cache.Get(url) 
	if ok {
		bytes = e
	}else {
		if err := getResBody(url, &bytes); err != nil { return nil, err }
		cache.Add(url, bytes)
	}

	var r LocAreaResult
	if err := json.Unmarshal(bytes, &r); err != nil {
		return nil, fmt.Errorf("Error while unmarshalling map: %v", err) 
	}
	return &r, nil
}

func getResBody(url string, outBytes *[]byte) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching map: %v", err)
	}
	if res.StatusCode > 299 || res.StatusCode < 200{
		return fmt.Errorf("Not found")
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error while reading body from url: %s : %v", url, err)
	}
	*outBytes = bytes
	return nil
}
