package configs

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	User         string `yaml:"user"`
	DbName       string `yaml:"dbname"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Sslmode      string `yaml:"sslmode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Timer        uint32 `yaml:"timer"`
	Db           string `yaml:"db_type"`
	ServerPort   string `yaml:"server_port"`
}

func ReadConfig() (*Config, error) {
	flag.Parse()
	var path string
	flag.StringVar(&path, "config_path", "config.yaml", "Путь к конфигу")

	dsnConfig := Config{}
	dsnFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dsnFile, &dsnConfig)
	if err != nil {
		return nil, err
	}

	return &dsnConfig, nil
}
