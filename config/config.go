package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config is the type structure of the config file for an environment
type Config struct {
	// Platform is the provisioning backend that will be used to create the
	// environment
	Platform Platform `yaml:"platform"`
	// Base contains configuration for the base image
	Base Base `yaml:"base"`
}

type Platform string

const (
	PlatformAWS = "aws"
)

// Base represents the steps necessary to build a base image for the current
// environment
type Base struct {
	// From is the ID of a base image to use as a starting point
	From string `yaml:"from"`
	// Run is an array of shell commands to be run on the base image
	Run []string `yaml:"run"`
}

// Load reads a config file
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
