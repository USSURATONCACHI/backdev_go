package main

import (
	"os"
	"io"
	"errors"
	"github.com/BurntSushi/toml"
)

type WrittenAppConfig struct {
	Secret string
	
	StartSyllables []string
	MiddleSyllables []string
	FinalSyllables []string
}

func GetConfigFromCli() (WrittenAppConfig, error) {
	// Get CLI argument
	if len(os.Args) != 2 {
		return WrittenAppConfig{}, errors.New("wrong amount of CLI arguments passed (must be 1)")
	}
	filePath := os.Args[1]
	
	// Open the file
	fsys := os.DirFS(".")
	file, err := fsys.Open(filePath)
	if err != nil {
		return WrittenAppConfig{}, errors.New("error opening config file: " + err.Error())
	}
	defer file.Close()

	// Read the file
	content, err := io.ReadAll(file)
	if err != nil {
		return WrittenAppConfig{}, errors.New("error reading config file: " + err.Error())
	}
	
	// Parse TOML
	var conf WrittenAppConfig
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		return WrittenAppConfig{}, errors.New("error parsing config file: " + err.Error())
	}

	return conf, nil
}