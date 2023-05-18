package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RedisConfig RedisConfig `yaml:"redis"`
}

func LoadFile(path string) (*Config, error) {
	c := Config{}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return &c, fmt.Errorf("failed to open %s file - %v", path, err)
	}
	if err = yaml.Unmarshal(buf, &c); err != nil {
		return &c, fmt.Errorf("failed to decode config file - %v", err)
	}
	return &c, nil
}
