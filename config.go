package main

import (
	"os"

	"github.com/kovetskiy/ko"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	UserID    string `yaml:"user_id",required:"true"`
	Secret    string `yaml:"secret",required:"true"`
	AccountID string `yaml:"account_id"`
	ProjectID string `yaml:"project_id"`
	Threads   int    `yaml:"threads"`

	path string
}

func NewConfig(path string) (Config, error) {
	config := Config{
		path: path,
	}

	err := ko.Load(path, &config, yaml.Unmarshal)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}

		return config, err
	}

	return config, nil
}
