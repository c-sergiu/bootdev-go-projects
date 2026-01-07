package pokeapi

import (
	"io"
	"fmt"
	"time"
	"net/http"
	"github.com/c-sergiu/pokego/internal/cache"
)

type PokeClient struct {
	Client *http.Client
	Cache *cache.Cache
}

func NewPokeClient(cacheInterval time.Duration) *PokeClient {
	return &PokeClient{
		Client: &http.Client{},
		Cache: cache.NewCache(cacheInterval),
	}
}

func (c *PokeClient) GetPokemon(name string) ([]byte, error) {
	if len(name) == 0 || name == "" {
		return nil, fmt.Errorf("invalid pokemon name")
	}
	base := "https://pokeapi.co/api/v2/pokemon/"
	url := base + name

	var bytes []byte
	if e, ok := c.Cache.Get(url); ok { bytes = e }else {
		if err := c.GetResBody(url, &bytes); err != nil { return nil, err }
		c.Cache.Add(url, bytes)
	}
	return bytes, nil
}

func (c *PokeClient) GetMap(name string) ([]byte, error) {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	url := baseUrl + name
	var bytes []byte
	if e, ok := c.Cache.Get(url); ok { bytes = e } else {
		if err := c.GetResBody(url, &bytes); err != nil {return nil, err}
		c.Cache.Add(url, bytes)
	}
	return bytes, nil
}

func (c *PokeClient) GetMaps(url string) ([]byte, error) {
	var bytes []byte
	if e, ok := c.Cache.Get(url); ok { bytes = e } else {
		if err := c.GetResBody(url, &bytes); err != nil { return nil, err }
		c.Cache.Add(url, bytes)
	}
	return bytes, nil
}

func (c *PokeClient) GetResBody(url string, outBytes *[]byte) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {return err}
	res, err := c.Client.Do(req)
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
