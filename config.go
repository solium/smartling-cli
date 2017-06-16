package main

import (
	"os"

	"github.com/gobwas/glob"
	"github.com/imdario/mergo"
	"github.com/kovetskiy/ko"
	hierr "github.com/reconquest/hierr-go"
	yaml "gopkg.in/yaml.v2"
)

type FileConfig struct {
	Pull struct {
		Format string
	}

	Push struct {
		Type       string
		Directives map[string]string
	}
}

type Config struct {
	UserID    string `yaml:"user_id",required:"true"`
	Secret    string `yaml:"secret",required:"true"`
	AccountID string `yaml:"account_id"`
	ProjectID string `yaml:"project_id"`
	Threads   int    `yaml:"threads"`

	Files map[string]FileConfig `yaml:"files"`

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

func (config *Config) GetFileConfig(path string) (FileConfig, error) {
	var match *FileConfig

	for key, candidate := range config.Files {
		pattern, err := glob.Compile(key, '/')
		if err != nil {
			return FileConfig{}, NewError(
				hierr.Errorf(
					err,
					`unable to compile pattern from config file (key "%s")`,
					key,
				),

				`File match pattern is malformed. Check out help for more `+
					`information on globbing patterns.`,
			)
		}

		if pattern.Match(path) {
			match = &candidate
		}
	}

	defaults := config.Files["default"]

	if match == nil {
		return defaults, nil
	}

	err := mergo.Merge(match, defaults)
	if err != nil {
		return FileConfig{}, NewError(
			hierr.Errorf(err, "unable to merge file config options"),
			`It's internal error. Consider reporting bug.`,
		)
	}

	return *match, nil
}
