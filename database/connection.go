package database

import (
	"database/sql"
	"fmt"
	_ "gopkg.in/goracle.v2"
	"log"
	"os"
)

const (
	atpUsername = "ATP_DEV_USERNAME"
	atpPassword = "ATP_DEV_PASSWORD"
	atpName     = "ATP_DEV_NAME"
)

var DB *sql.DB

func Initialize() {
	config := dbConfig()

	connectionString := fmt.Sprintf("%s/%s@%s", config[atpUsername], config[atpPassword], config[atpName])
	var err error
	DB, err = sql.Open("goracle", connectionString)
	if err != nil {
		log.Fatal(err.Error())
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected Oracle DB!")

}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	username, ok := os.LookupEnv(atpUsername)
	if !ok {
		panic(atpUsername + " missing from .env file")
	}

	password, ok := os.LookupEnv(atpPassword)
	if !ok {
		panic(atpPassword + " missing from .env file")
	}

	name, ok := os.LookupEnv(atpName)
	if !ok {
		panic(atpName + " missing from .env file")
	}

	conf[atpUsername] = username
	conf[atpPassword] = password
	conf[atpName] = name
	return conf
}
