package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"path/filepath"
	"io/ioutil"
	"github.com/labstack/gommon/log"
)

type DbKitConfig struct {
	Db struct {
		Type     string
		Host     string
		Port     int
		Username string
		Password string
		Database []string
	}

	CronExpression string

	Storage struct {
		Type    string
		Account string
	}
}

func main() {
	var config *DbKitConfig
	configPath := filepath.Join(".", "config.yml")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(config)
}
