package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"time"
)

type Config struct {
	PathToData string     `yaml:"pathToData"`
	Server     HttpServer `yaml:"server"`
}

type HttpServer struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	ShutdownTime time.Duration `yaml:"shutdownTime"`
}

func MustNew(cfgString string) Config {
	file, err := os.Open(cfgString)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("config file is not set")
		}
		log.Fatal("error while opening config file", err)
	}

	cfgData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("error while reading file", err)
	}

	config := Config{}
	if err = yaml.Unmarshal(cfgData, &config); err != nil {
		log.Fatal("cannot read config:", err)
	}

	return config
}
