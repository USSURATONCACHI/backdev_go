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

func modelSmtpInfoFromConfig(Smtp config.Smtp) model.SmtpInfo {
	return model.SmtpInfo{
		Host: Smtp.Host,
		Port: Smtp.Port,
		User: Smtp.User,

		Password:  Smtp.Password,
		FromEmail: Smtp.FromEmail,

		MockUserEmail: Smtp.MockUserEmail,
	}
}

func modelSyllablesFromConfig(Syllables config.Syllables) model.Syllables {
	return model.Syllables{
		Start: Syllables.Start,
		Middle: Syllables.Middle,
		Final: Syllables.Final,
	}
}

func ModelFromConfig(config config.App) (model.Model, func(), error) {
	// Convert base info to model
	mdl := model.Model {
		Secret: sha512.Sum512([]byte(config.Secret)),
		Syllables: modelSyllablesFromConfig(config.Syllables),
		SmtpInfo: modelSmtpInfoFromConfig(config.Smtp),
	}
	deferFunc := func() {}

	// Choose database type
	if config.DbType == "in_memory" {
		// IN-MEMORY DATABASE
		mdl.Database = db_io.InMemoryDatabaseNew()
	} else if config.DbType == "postgresql" {
		// POSTGRESQL DATABASE
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
		// INVALID DATABASE
		return mdl, deferFunc, errors.New("unknown database type. Only supported are: 'in_memory', 'postgresql'")
	}

	return mdl, deferFunc, nil
}