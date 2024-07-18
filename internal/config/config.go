package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env"`
	DSN      string `yaml:"dsn"`
	GRPC     `yaml:"grpc"`
	TokenTTL time.Duration `yaml:"token_ttl"`
}

type GRPC struct {
	Port int `yaml:"port"`
}

func MustLoad() Config {
	p := fetchPath()
	c := Config{}
	if err := cleanenv.ReadConfig(p, &c); err != nil {
		panic(err.Error())
	}

	return c
}

func fetchPath() string {
	var path string
	flag.StringVar(&path, "cfg", "", "path to service config")
	flag.Parse()
	if path != "" {
		return path
	}
	path = os.Getenv("SSO_CONFIG_PATH")
	if path != "" {
		return path
	}

	panic("path to config file is not defined")
}
