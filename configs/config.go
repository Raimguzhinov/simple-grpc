package configs

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	EventServer struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
	} `yaml:"server"`

	CaldavServer struct {
		Url      string `yaml:"url"`
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"caldav"`

	RabbitmqServer struct {
		Url         string `yaml:"url"`
		Login       string `yaml:"login"`
		Password    string `yaml:"password"`
		Exchange    string `yaml:"exchange"`
		ReconnDelay uint   `yaml:"reconnection_delay"`
	} `yaml:"rabbitmq"`
}

func New(configPath string) (*Config, error) {
	config := &Config{}
	if err := validateConfigPath(configPath); err != nil {
		return nil, err
	}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
