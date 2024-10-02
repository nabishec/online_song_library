package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string `yaml:"env" env-default:"local" env-required:"true"`
	DBDataSourceName `yaml:"db_data_source_name"`
	HTTPServer       `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type DBDataSourceName struct {
	Protocol string `yaml:"protocol" env-default:"postgres"`
	UserName string `yaml:"user_name" env-default:"postgres"`
	Password string `yaml:"password" env-default:"secret"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"dbname" env-default:"postgres"`
	Options  string `yaml:"options" env-default:"sslmode=disable"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Println("CONFIG_PATH isn't set")
		configPath = "./config/local.yml"
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
