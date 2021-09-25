package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	var c string
	flag.StringVar(&c, "c", "app.yml", "config file path")
	flag.Parse()

	configOption, err := GetConfigOption(c)
	if err != nil {
		return
	}

	apple := NewApple(&AppleOption{
		ConfigOption: configOption,
	})

	pid := os.Getpid()
	log.Printf("[I] pid:%v", pid)

	apple.Serve()
}

func GetConfigOption(c string) (*ConfigOption, error) {
	configOption := &ConfigOption{}

	data, err := ioutil.ReadFile(c)
	if err != nil {
		log.Printf("[E] read config failed. file:%v, err:%v", c, err)
		return nil, err
	}

	err = yaml.Unmarshal(data, configOption)
	if err != nil {
		log.Printf("[E] unmarshal config failed. err:%v", err)
		return nil, err
	}

	return configOption, nil
}
