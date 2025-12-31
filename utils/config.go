package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Upload struct {
		Directory string `yaml:"directory"`
		MaxSize   int    `yaml:"max_size"`
	} `yaml:"upload"`
	Deploy struct {
		Backend struct {
			ScriptPath string `yaml:"script_path"`
			WorkingDir string `yaml:"working_dir"`
		} `yaml:"backend"`
		Frontend struct {
			ScriptPath string `yaml:"script_path"`
			WorkingDir string `yaml:"working_dir"`
		} `yaml:"frontend"`
	} `yaml:"deploy"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
