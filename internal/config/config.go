package config

import "time"

type Config struct {
	Env          string  `yaml:"env"`
	Storage_path string  `yaml:"storage_path"`
	Storage      Storage `yaml:"storage"`
	Cache        Cache   `yaml:"cache"`
	Grpc         Grpc    `yaml:"grpc"`
}

type Storage struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
}

type Cache struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type Grpc struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
