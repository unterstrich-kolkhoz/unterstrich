package config

import (
	"errors"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

// Config is our global configuration object
type Config struct {
	Port              string
	Staticdir         string `toml:"static_dir"`
	Uploader          string
	ImgThumbnailer    string `toml:"img_thumbnailer"`
	VidThumbnailer    string `toml:"vid_thumbnailer"`
	ShaderThumbnailer string `toml:"shader_thumbnailer"`
	TemplateDir       string `toml:"template_dir"`
	SQLDialect        string `toml:"sql_dialect"`
	SQLName           string `toml:"db_name"`
	StripeKey         string `toml:"stripe_key"`
	Region            string `toml:"region"`
	URL               string `toml:"url"`
	HostedZoneID      string `toml:"hosted_zone_id"`
	S3HostedZoneID    string `toml:"s3_hosted_zone_id"`
	S3URL             string `toml:"s3_url"`
}

// ReadConfig is the main entry point for configuration
// parsing
func ReadConfig(configFile string) (*Config, error) {
	configFiles := configOptions(configFile)
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

func configOptions(filename string) []string {
	return []string{
		strings.Join([]string{"/", filename}, ""),
		strings.Join([]string{"./", filename}, ""),
		strings.Replace(filename, ".conf", ".local.conf", 1),
	}
}
