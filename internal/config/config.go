package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".BatchBaseconfig.json"

type Config struct {
	Dburl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homeDir, configFileName)
	return path, nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) SetUser(username string) error {

	cfg.CurrentUserName = username
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
