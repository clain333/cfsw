package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

var C Config

type Config struct {
	Host string `yaml:"host"`
	Auth string `yaml:"auth"`
	Cidr string `yaml:"cidr"`
}

func LoadYaml(filename string) error {

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &C)
	if err != nil {
		return err
	}

	return nil
}
