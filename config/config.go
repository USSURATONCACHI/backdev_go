package config

import (
	"errors"
	"io"
	"os"
	"strconv"

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
}

type App struct {
	Secret string
	DbType string
	ListenIp string

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

	UpdateConfigFromEnv(&conf)

	return conf, nil
}

func UpdateConfigFromEnv(conf *App) {
	if val, is_set := os.LookupEnv("BACKDEV_SECRET"); is_set {
		conf.Secret = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_DB_TYPE"); is_set {
		conf.DbType = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_LISTEN_IP"); is_set {
		conf.ListenIp = val
	}

	// PSQL
	if val, is_set := os.LookupEnv("BACKDEV_PSQL_HOST"); is_set {
		conf.Postgresql.Host = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_PSQL_PORT"); is_set {
		parsed, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			panic("Failed to parse BACKDEV_PSQL_PORT to an int16")
		}
		conf.Postgresql.Port = int16(parsed)
	}
	if val, is_set := os.LookupEnv("BACKDEV_PSQL_USER"); is_set {
		conf.Postgresql.User = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_PSQL_PASSWORD"); is_set {
		conf.Postgresql.Password = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_PSQL_DBNAME"); is_set {
		conf.Postgresql.DbName = val
	}

	// SMTP
	if val, is_set := os.LookupEnv("BACKDEV_SMTP_HOST"); is_set {
		conf.Smtp.Host = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_SMTP_PORT"); is_set {
		parsed, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			panic("Failed to parse BACKDEV_SMTP_PORT to an int16")
		}
		conf.Smtp.Port = int16(parsed)
	}
	if val, is_set := os.LookupEnv("BACKDEV_SMTP_USER"); is_set {
		conf.Smtp.User = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_SMTP_PASSWORD"); is_set {
		conf.Smtp.Password = val
	}
	if val, is_set := os.LookupEnv("BACKDEV_SMTP_FROM_EMAIL"); is_set {
		conf.Smtp.FromEmail = val
	}
}