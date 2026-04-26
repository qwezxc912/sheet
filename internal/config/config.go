package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env" env-default:"local" env-required:"true"`
	Server `yaml:"server"`
}

type Server struct {
	Addres      string        `yaml:"addres" env-default: "localhost:8084"`
	Timeout     time.Duration `yaml:"timeout" env_default: "4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default: "60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH isn't set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s\n", err)
	}

	return &cfg
}
