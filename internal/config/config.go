package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type LicenseConfig struct {
	LisDir string `yaml:"license_dir"`
}

func ReadConf(file string) (*LicenseConfig, error) {
	conf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg LicenseConfig
	err = yaml.Unmarshal(conf, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
