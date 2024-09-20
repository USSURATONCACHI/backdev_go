package main

import (
	"os"
	"io"
	"errors"
	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	Secret string
}

func GetConfigFromCli() (AppConfig, error) {
	// Get CLI argument
	if len(os.Args) != 2 {
		return AppConfig{}, errors.New("Wrong amount of CLI arguments passed (must be 1)")
	}
	filePath := os.Args[1]
	
	// Open the file
	fsys := os.DirFS(".")
	file, err := fsys.Open(filePath)
	if err != nil {
		return AppConfig{}, errors.New("Error opening config file: " + err.Error())
	}
	defer file.Close()

	// Read the file
	content, err := io.ReadAll(file)
	if err != nil {
		return AppConfig{}, errors.New("Error reading config file: " + err.Error())
	}
	
	// Parse TOML
	var conf AppConfig
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		return AppConfig{}, errors.New("Error parsing config file: " + err.Error())
	}

	return conf, nil
}