package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"addr"`
}
type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HttpServer  `yaml:"http_server"`
}

func MustLoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "provide your config path here dude")

		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("The config path is not set yet")
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exist %s", configPath)
		}

	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read the config file %s", err.Error())
	}

	return &cfg
}
