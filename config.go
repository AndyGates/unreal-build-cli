package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type OptionSet struct {
	Options  []string
	Defaults []int
}

type Config struct {
	RunUATPath           string
	ClientOptions        OptionSet
	ServerOptions        OptionSet
	ConfigurationOptions OptionSet
	StepOptions          OptionSet
}

func GetConfig() Config {

	const configFile = "unreal-build-cli.config.json"

	config := CreateDefaultConfig()

	if CheckFileExists(configFile) {
		log.Printf("Loading config from file: %s", configFile)

		file, _ := os.ReadFile(configFile)
		json.Unmarshal(file, &config)
	}

	fmt.Println(config)

	return config
}

func CreateDefaultConfig() Config {
	return Config{
		RunUATPath: "..\\UnrealEngine\\Engine\\Build\\BatchFiles\\RunUAT.bat",
		ClientOptions: OptionSet{
			Options:  []string{"Win64", "PS5", "XSX"},
			Defaults: []int{0},
		},
		ServerOptions: OptionSet{
			Options:  []string{"Win64", "Linux"},
			Defaults: []int{},
		},
		ConfigurationOptions: OptionSet{
			Options:  []string{"Development", "Shipping", "Test", "Debug"},
			Defaults: []int{0},
		},
		StepOptions: OptionSet{
			Options:  []string{"Build", "Cook", "Pak", "Stage"},
			Defaults: []int{0, 3},
		},
	}
}
