package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env            string        `yaml:"env"`
	MigrationsPath string        `yaml:"migrations_path"`
	Storage        Storage       `yaml:"storage"`
	Cache          Cache         `yaml:"cache"`
	Grpc           Grpc          `yaml:"grpc"`
	TokenTTL       time.Duration `yaml:"tokenTTL"`
}

type Storage struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"dbname"`
	Username string `yaml:"username"`
}

type Cache struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DB     int    `yaml:"db"`
}

type Grpc struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	return MustLoadFromPath(os.Getenv("CONFIG_PATH"))
}

func MustLoadFromPath(path string) *Config {
	data, err := os.ReadFile(path)
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
