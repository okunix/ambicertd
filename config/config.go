package config

import (
	"io"
	"os"

	"go.yaml.in/yaml/v4"
)

type TLSConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type ServerConfig struct {
	Addr string     `yaml:"addr"`
	TLS  *TLSConfig `yaml:"tls,omitempty"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	SQLite SQLiteConfig `yaml:"sqlite"`
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr: "0.0.0.0:6767",
			TLS:  nil,
		},
		SQLite: SQLiteConfig{
			Path: "ambicert.db",
		},
	}
}

func Parse(r io.Reader) (Config, error) {
	cfg := defaultConfig()
	err := yaml.NewDecoder(r).Decode(&cfg)
	return cfg, err
}

func ReadFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return defaultConfig(), err
	}
	defer file.Close()
	return Parse(file)
}
