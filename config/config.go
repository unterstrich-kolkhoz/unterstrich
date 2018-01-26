package config

import (
	"errors"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

type Config struct {
	Port string
}

func ReadConfig(configFile string) (*Config, error) {
	configFiles := ConfigOptions(configFile)
	config := &Config{}
	hasConfig := false
	var confError error

	for _, filename := range configFiles {
		tmp := &Config{}
		_, err := toml.DecodeFile(filename, tmp)
		if err != nil {
			continue
		} else {
			log.Println("Using config file:", filename)
			hasConfig = true
			// Merge configs
			if err := mergo.Merge(config, tmp); err != nil {
				return nil, err
			}
		}
	}

	if !hasConfig {
		confError = errors.New("Could not load any config file")
	}

	return config, confError
}

func ConfigOptions(filename string) []string {
	return []string{
		strings.Join([]string{"/", filename}, ""),
		strings.Join([]string{"./", filename}, ""),
		strings.Replace(filename, ".conf", ".local.conf", 1),
	}
}
