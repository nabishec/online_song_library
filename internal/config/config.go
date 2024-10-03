package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	err := godotenv.Load("configuration.env")
	if err != nil {
		log.Fatal("Mistake of download .env:", slerr.Err(err))
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Println("CONFIG_PATH isn't set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file isn't exists:", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config:", err)
	}
	return &cfg
}
