package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"math/rand"
	"net/http"
	"encoding/json"
	"io"
	"github.com/c-sergiu/pokego/internal/cache"
	"github.com/c-sergiu/pokego/internal/repl"
)

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
	Height         int    `json:"height"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}


type Pokedex struct {
	Dex map[string]Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		Dex: make(map[string]Pokemon),
	}
}

func (p *Pokedex) Add(pokemon Pokemon) {
	p.Dex[pokemon.Name] = pokemon
}

func (p Pokedex) Print() {
	fmt.Println("Your Pokedex:")
	for _, p := range p.Dex {
		fmt.Printf(" - %s\n", p.Name)
	}
}

type Encounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type LocArea struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type LocAreaResult struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Prev *string `json:"previous"`
	Results []LocArea `json:"results"`
}

func (l LocAreaResult) DisplayResults() {
	for _, loc := range l.Results {
		fmt.Println(loc.Name)
	}
}

type NavConfig struct {
	Next *string
	Prev *string
}

func NewConfig(url string) *NavConfig {
	return &NavConfig{
		Next: &url,
		Prev: nil,
	}
}

func (c *NavConfig) Update(prev *string, next *string) {
	c.Next = next
	c.Prev = prev
}

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
