package main

import (
	"events/task/external/client"
	"events/task/external/logger"
	"events/task/internal/event/repository"
	"flag"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Host string
	Port int
}

func main() {
	var pathToCfgFile string
	flag.StringVar(&pathToCfgFile, "cfg", "cfg.yaml", "set path to config file")

	file, err := os.Open(pathToCfgFile)
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

	lg := logger.New()

	fileClient, err := client.GetFSClient(lg)
	if err != nil {
		lg.Error(err.Error())
	}

	repository, err := repository.NewRepository(fileClient, lg)
	if err != nil {
		lg.Error(err.Error())
	}

	router := http.NewServeMux()

}
