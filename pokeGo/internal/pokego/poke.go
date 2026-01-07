package pokego

import (
	"fmt"
	"time"
	"github.com/c-sergiu/pokego/internal/pokeapi"
)

type Context struct {
	Client *pokeapi.PokeClient
	Nav *NavConfig
	Dex map[string]Pokemon
}

func NewContext() *Context {
	return &Context{
		Client: pokeapi.NewPokeClient(60 * time.Second),
		Nav: NewNavConfig("https://pokeapi.co/api/v2/location-area/"),
		Dex: make(map[string]Pokemon),
	}
}

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

func NewNavConfig(url string) *NavConfig {
	return &NavConfig{
		Next: &url,
		Prev: nil,
	}
}

func (c *NavConfig) Update(prev *string, next *string) {
	c.Next = next
	c.Prev = prev
}
