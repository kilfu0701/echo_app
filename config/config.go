package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	DBName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Timeout  int    `yaml:"timeout"`
}

type Config struct {
	EnvName string `yaml:"env_name"`

	Database DBConfig `yaml:"db"`
}

func Load(configFilePath string) *Config {
	absPath, _ := filepath.Abs(configFilePath)
	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Printf("yaml read err   #%v ", err)
	}

	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &c
}
