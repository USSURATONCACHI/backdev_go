package config

import (
	"os"
	"io"
	"errors"
	"github.com/BurntSushi/toml"
)

type Syllables struct {
	Start []string
	Middle []string
	Final []string
}

type Database struct {
	Host string
	Port int16
	User string
	Password string
	DbName string
}

type Smtp struct {
	Host string
	Port int16
	User string
	Password string

	FromEmail string
	MockUserEmail string
}

type App struct {
	Secret string
	DbType string

	Syllables Syllables
	Postgresql Database
	
	Smtp Smtp
}

func GetConfigFromCli() (App, error) {
	// Get CLI argument
	if len(os.Args) != 2 {
		return App{}, errors.New("wrong amount of CLI arguments passed (must be 1)")
	}
	filePath := os.Args[1]
	
	// Open the file
	fsys := os.DirFS(".")
	file, err := fsys.Open(filePath)
	if err != nil {
		return App{}, errors.New("error opening config file: " + err.Error())
	}
	defer file.Close()

	// Read the file
	content, err := io.ReadAll(file)
	if err != nil {
		return App{}, errors.New("error reading config file: " + err.Error())
	}
	
	// Parse TOML
	var conf App
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		return App{}, errors.New("error parsing config file: " + err.Error())
	}

	return conf, nil
}