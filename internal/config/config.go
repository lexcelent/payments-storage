package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	StoragePath string     `yaml:"storage_path"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

func MustLoad() *Config {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
