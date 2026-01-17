package config

import (
	"encoding/json"
	"os"
	"fmt"
)

const configFileName = "/.gatorconfig.json"

type Config struct{
	DbUrl string   `json:"db_url"`
	CurUser string `json:"current_user_name"`
}

func Read() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err !=  nil {return nil, err}
	defer f.Close()
	
	decoder := json.NewDecoder(f)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurUser = user
	if err := write(cfg); err != nil {
		return err
	}
	return nil
}

func (cfg Config) Display() string {
	return fmt.Sprintf("User: %s\nDb: %s\n", cfg.CurUser, cfg.DbUrl)
}

func write(cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err !=  nil {return err}
	defer f.Close()
	
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home + configFileName
	return path, nil
}


