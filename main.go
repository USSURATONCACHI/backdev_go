package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"crypto/sha512"

	"backdev_go/config"
	"backdev_go/db_io"
	"backdev_go/model"
)


func main() {
	// Read config
	config, err := config.GetConfigFromCli()
	if err != nil {
		fmt.Println("Failed to parse config")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create model
	mdl, deferFunc, err := ModelFromConfig(config)
	if err != nil {
		fmt.Println("Failed to create model")
		fmt.Println(err.Error())
		os.Exit(2)
	}
	defer deferFunc()

	fmt.Println("Your base64 of secret: ", base64.StdEncoding.EncodeToString(mdl.Secret[:]))
	

	// Run server
	server := CreateServer(mdl)
	server.Run("localhost:8080")
}

func ModelFromConfig(config config.App) (model.Model, func(), error) {
	mdl := model.Model {
		Secret: sha512.Sum512([]byte(config.Secret)),
		Syllables: model.Syllables {
			Start: config.Syllables.Start,
			Middle: config.Syllables.Middle,
			Final: config.Syllables.Final,
		},
	}
	deferFunc := func() {}

	if config.DbType == "in_memory" {
		mdl.Database = db_io.InMemoryDatabaseNew()
	} else if config.DbType == "postgresql" {
		requests := db_io.PostgresqlRequestsDefault()
		options := db_io.DatabaseParams {
			Host: config.Postgresql.Host,
			Port: config.Postgresql.Port,
			User: config.Postgresql.User,
			Password: config.Postgresql.Password,
			DbName: config.Postgresql.DbName,
		}

		db, err := db_io.PostgresqlDatabaseNew(options, requests)
		if err != nil {
			return mdl, deferFunc, err
		}
		deferFunc = func() { db.Close() }
		mdl.Database = db
	} else {
		return mdl, deferFunc, errors.New("unknown database type. Only supported are: 'in_memory', 'postgresql'")
	}

	return mdl, deferFunc, nil
}