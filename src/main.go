package main

import (
	"backend-sample/database"
	"encoding/json"
	"log"
)

var dbConfig string = `{
	"host": "",
	"database": "",
	"user": "",
	"password": "",
	"port": 3306,
	"maxLifetime": 60,
	"maxOpenConns": 5,
	"maxIdleConns": 5
}`

func main() {
	parseDbConfig()

	user := database.User{Name: "Lorem", Email: "lorem@mail.com", Password: "password"}
	if !database.CreateUser(&user) {
		log.Println("Coudl not create user")
	}
}

func parseDbConfig() {
	var config database.DatabaseConfiguration
	err := json.Unmarshal([]byte(dbConfig), &config)

	if err != nil {
		panic("Could not parse database config.")
	}

	database.Configuration = config
}
