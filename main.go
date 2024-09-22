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
	"backdev_go/smtp_io"
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
	fmt.Printf("Running on IP '%s'\n", config.ListenIp)
	server.Run(config.ListenIp)
}


func ModelFromConfig(config config.App) (model.Model, func(), error) {
	// Convert base info to model
	mdl := model.Model {
		Secret: sha512.Sum512([]byte(config.Secret)),
		Syllables: modelSyllablesFromConfig(config.Syllables),
		SmtpClient: smtpClientFromConfig(config.Smtp),
	}
	
	db, deferFunc, error := databaseFromConfig(config)
	mdl.Database = db

	return mdl, deferFunc, error
}

func databaseFromConfig(config config.App) (db_io.Database, func(), error) {
	if config.DbType == "in_memory" {
		// IN-MEMORY DATABASE
		return db_io.InMemoryDatabaseNew(), func(){}, nil;

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
			return nil, func(){}, err
		}
		deferFunc := func() { db.Close() }
		
		return db, deferFunc, nil

	} else {
		// INVALID DATABASE
		return nil, func(){}, errors.New("unknown database type. Only supported are: 'in_memory', 'postgresql'")
	}
}

func smtpClientFromConfig(Smtp config.Smtp) smtp_io.SmtpClient {
	return &smtp_io.PlainAuthClient {
		Host: Smtp.Host,
		Port: Smtp.Port,
		User: Smtp.User,

		Password:  Smtp.Password,
		FromEmail: Smtp.FromEmail,
	}
}

func modelSyllablesFromConfig(Syllables config.Syllables) model.Syllables {
	return model.Syllables{
		Start: Syllables.Start,
		Middle: Syllables.Middle,
		Final: Syllables.Final,
	}
}