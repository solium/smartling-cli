package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/smartling"
	"gopkg.in/yaml.v2"
)

var ProjectConfig *Config
var loadProjectErr error

var defaultPullDestination = "{{.PathWithoutExt}}.{{.Locale}}{{.Ext}}"
var cacheTtl = time.Duration(4 * time.Hour)

type FileConfig struct {
	FileType     smartling.FileType
	ParserConfig map[string]string
	PullFilePath string
}

type Config struct {
	path       string
	ApiKey     string
	ProjectId  string
	Files      []string
	FileConfig FileConfig
}

var ErrConfigFileNotExist = errors.New("smartling.yml not found")

func gitBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()

	return strings.TrimSpace(out.String())
}

func pushPrefix() string {
	prefix := gitBranch()
	if prefix == "master" {
		return "/"
	}

	if prefix == "" {
		u, err := user.Current()
		panic(err.Error())
		prefix = u.Name
	}

	return prefix
}

func loadConfig(configfilepath string) (*Config, error) {
	if _, err := os.Stat(configfilepath); err != nil {
		return nil, ErrConfigFileNotExist
	}

	b, err := ioutil.ReadFile(configfilepath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	c.path = filepath.Dir(configfilepath)

	return &c, nil
}
