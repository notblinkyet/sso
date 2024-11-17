package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env     string  `yaml:"env"`
	Storage Storage `yaml:"storage"`
	Cache   Cache   `yaml:"cache"`
	Grpc    Grpc    `yaml:"grpc"`
}

type Storage struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"dbname"`
	Username string `yaml:"username"`
}

type Cache struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
}

type Grpc struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	data, err := os.ReadFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
