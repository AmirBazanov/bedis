package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Logger Logger
	Server Server
}

type Logger struct {
	Level   string `yaml:"level"`
	Service string `yaml:"service"`
	Logfile string `yaml:"logfile"`
}

type Server struct {
	Port    string `yaml:"port"`
	Address string `yaml:"address"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	return MustLoadPath(configPath)

}

func MustLoadPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file not found: " + path)
	}
	var config Config
	if cleanenv.ReadConfig(path, &config) != nil {
		panic("Config file not found: " + path)
	}
	return &config
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
