package config

import (
	"flag"
)

type Config struct {
	Project     *string
	Secret      *string
	Namespace   *string
	ClusterMode *bool
	ConfigFile  *string
	Terraformrc *bool
	ShowVersion *bool
	Context     *string
}

func NewConfig() *Config {
	return &Config{
		Project:     flag.String("project", "none", "project name for tfvars"),
		Secret:      flag.String("secret", "test", "Secret name for tfvars"),
		Namespace:   flag.String("n", "default", "Namespace for secret"),
		ClusterMode: flag.Bool("cluster", true, "Set -cluster=false for local test"),
		ConfigFile:  flag.String("config", "~/.kube/config", "Config file path. Default ~/.kube/config"),
		Terraformrc: flag.Bool("terraformrc", false, "Choose terraformrc mirror"),
		ShowVersion: flag.Bool("version", false, "Print version"),
		Context:     flag.String("context", "pulse/agents:dev-test", "Context for kubeconfig"),
	}
}

func (c *Config) Parse() {
	flag.Parse()
}
