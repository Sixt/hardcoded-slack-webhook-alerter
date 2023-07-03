package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ConfigPath is the default path to the config file
const ConfigPath = "config/config.yaml"

// Config contains the configuration for the scanner
type Config struct {
	ChannelPattern string   `yaml:"channel_pattern"`  //used to extract the slack channel from the file
	Directories    []string `yaml:"directories,flow"` // Directories is the list of paths to scan for hardcoded credentials
	GithubOrg      string   `yaml:"github_org"`
	DryRun         bool     `yaml:"dry_run"`
}

// LoadConfig loads the configuration from the config.yaml file
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, err
	}

	return c, nil
}
