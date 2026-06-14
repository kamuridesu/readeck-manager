package config

import (
	"errors"
	"os"
)

type Config struct {
	ReadeckURL   string `yaml:"readeck_url"`
	ReadeckToken string `yaml:"readeck_token"`
	OpenAIUrl    string `yaml:"openai_url"`
	OpenAIKey    string `yaml:"openai_key"`
	OpenAIModel  string `yaml:"openai_model"`
}

func Load() (*Config, error) {
	cfg := &Config{
		ReadeckURL:   os.Getenv("READECK_URL"),
		ReadeckToken: os.Getenv("READECK_TOKEN"),
		OpenAIUrl:    os.Getenv("OPENAI_URL"),
		OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:  os.Getenv("OPENAI_MODEL"),
	}

	if cfg.ReadeckURL == "" {
		return nil, errors.New("READECK_URL environment variable is required")
	}
	if cfg.ReadeckToken == "" {
		return nil, errors.New("READECK_TOKEN environment variable is required")
	}
	if cfg.OpenAIKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}
	if cfg.OpenAIModel == "" {
		cfg.OpenAIModel = "gpt-4o-mini"
	}

	return cfg, nil
}

