package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
	Multiple   `yaml:"multiple"`
}

type HttpServer struct {
	Port          string `yaml:"port"`
	Host          string `yaml:"host"`
	Timeout       string `yaml:"timeout"`
	IdleTimeout   string `yaml:"idle_timeout"`
	HeaderTimeout string `yaml:"header_timeout"`
}
type Postgres struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
}

type Multiple struct {
	MaxRacer int `yaml:"max_racers"`
	Timer    int `yaml:"timer"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
		return nil, fmt.Errorf("unable read from config: %w", err)
	}
	return &cfg, nil
}
