package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Slack struct {
		Token string
	}
	Icinga2 struct {
		Host     string
		Username string
		Password string
	}
}

func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
