package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
	Multiple   `yaml:"multiple"`
	AWS        AWS
	Jwt        JWT
	Admin      Admin
}

type HttpServer struct {
	Port          string `yaml:"port"`
	Host          string `yaml:"host"`
	Timeout       string `yaml:"timeout"`
	IdleTimeout   string `yaml:"idle_timeout"`
	HeaderTimeout string `yaml:"header_timeout"`
}
type Postgres struct {
	Port     string `env:"POSTGRES_PORT"`
	Host     string `env:"POSTGRES_HOST"`
	Database string `env:"POSTGRES_DB"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	PgUrl    string `env:"PG_URL"`
}

type Multiple struct {
	MaxRacer int `yaml:"max_racers"`
	Timer    int `yaml:"timer"`
}

type AWS struct {
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	Region          string `env:"AWS_REGION"`
	BucketName      string `env:"AWS_BUCKET_NAME"`
}

type JWT struct {
	SecretKey string `env:"JWT_SECRET_KEY"`
}

type Admin struct {
	Username string `env:"ADMIN_USERNAME"`
	Password string `env:"ADMIN_PASSWORD"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
		return nil, fmt.Errorf("unable read from config: %w", err)
	}
	return &cfg, nil
}
