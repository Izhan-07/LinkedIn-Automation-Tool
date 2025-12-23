package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Browser  BrowserConfig  `yaml:"browser"`
	LinkedIn LinkedInConfig `yaml:"linkedin"`
	Database DatabaseConfig `yaml:"database"`
	Limits   LimitsConfig   `yaml:"limits"`
}

type AppConfig struct {
	Debug bool `yaml:"debug"`
}

type BrowserConfig struct {
	Headless bool   `yaml:"headless"`
	Stealth  bool   `yaml:"stealth"`
	UserData string `yaml:"user_data"` // Path to Chrome user data dir
}

type LinkedInConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type LimitsConfig struct {
	DailyConnections int `yaml:"daily_connections"`
	DailyMessages    int `yaml:"daily_messages"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Env override for credentials
	if user := os.Getenv("LINKEDIN_USERNAME"); user != "" {
		cfg.LinkedIn.Username = user
	}
	if pass := os.Getenv("LINKEDIN_PASSWORD"); pass != "" {
		cfg.LinkedIn.Password = pass
	}

	return &cfg, nil
}
