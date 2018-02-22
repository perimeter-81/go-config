package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/neliseev/go-config"
)

type params struct {
	Home         string        `json:"home" env:"HOME" envDefault:"/home/some"`
	Port         int           `yaml:"port" env:"PORT" envDefault:"3000"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":" envDefault:"8.8.8.8:8.8.4.4"`
	Duration     time.Duration `env:"DURATION"`
	JSONSubKey   struct {
		SomeKey string `json:"some_key"`
	} `json:"json_sub_key"`
	YAMLSubKey struct {
		SomeKey string `yaml:"some_key"`
	} `yaml:"yaml_sub_key"`
}

func main() {
	var err error

	cfg := new(params)
	if err = config.Parse(nil, cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ENV only: %+v\n", cfg)

	jsData, err := ioutil.ReadFile("./_test/test.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JSON:\n%s\n", jsData)

	if err = config.Parse(jsData, cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON + ENV: %+v\n", cfg)

	ymlData, err := ioutil.ReadFile("./_test/test.yml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("YML:\n%s\n", ymlData)

	if err = config.Parse(ymlData, cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON + YAML: %+v\n", cfg)
}
