package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)


type HTTPServer struct {
	Addr string `yaml:"address" env-default:"localhost:8000"`
}

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
    DatabaseURL string `yaml:"DatabaseURL" env-required:"true"`
    DatabaseName string `yaml:"DatabaseName" env-required:"true"`
	JwtSecret    string `yaml:"JwtSecret"`
	HTTPServer `yaml:"http_server"`
}


func MustLoad() *Config{
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			panic("config path is required")
		}

	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	return &cfg
}