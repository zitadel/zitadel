package config

import (
	"github.com/caos/utils/config/yaml"
	"github.com/caos/utils/logging"
	log "github.com/caos/utils/logging/config"
)

type Config struct {
	Port      string
	StaticDir string
	Log       *log.Log
}

func MustReadConfig(path string) *Config {
	config, err := ReadConfig(path)
	logging.Log("CONFI-W5q2O").OnError(err).Panic("unable to read config")

	return config
}

func ReadConfig(path string) (*Config, error) {
	config := new(Config)
	err := yaml.ReadConfig(config, path)
	//TODO: runtime config from vault or something like that
	//TODO: process to autoamtically scrap runtime config changes
	return config, err
}
