package main

import (
	"backend-sample/apis"
	"backend-sample/common"
	"backend-sample/database"
	"backend-sample/middlewares"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	errInvalidConfigFile = errors.New("invalid configuration file")
)

var mysqldb database.MySqlDatabaseService

func main() {
	readDbConfig()
	apis.Initialize(mysqldb)

	key, err := common.GenerateAESKey(32)

	if err != nil {
		log.Fatalf("Could not generate AES")
	}

	fmt.Println(fmt.Sprintf("AES key: %s", common.EncodeBase64(key)))

	router := gin.Default()
	router.Use(middlewares.MiddlewareHandler)
	router.Use(middlewares.ResponseFormatMiddleware)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/users", apis.GetUser)
	router.POST("/users", apis.AddUser)
	router.DELETE("/users/:userId", apis.DeleteUser)
	router.PUT("/users/:userId", apis.UpdateUser)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

func readDbConfig() {
	viper.SetConfigName("db")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../configs")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		panic(errInvalidConfigFile)
	}

	// err := yaml.Unmarshal([]byte(dbConfig), &config)

	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(errCouldNotParseConfigFile)
	// }

	mysqldb.Configuration = database.DatabaseConfiguration{
		Host:         viper.GetString("database.Host"),
		Database:     viper.GetString("database.Database"),
		User:         viper.GetString("database.User"),
		Password:     viper.GetString("database.Password"),
		Port:         viper.GetInt("database.Port"),
		MaxLifetime:  viper.GetInt("database.MaxLifetime"),
		MaxOpenConns: viper.GetInt("database.MaxOpenConns"),
		MaxIdleConns: viper.GetInt("database.MaxIdleConns"),
	}

	var keys []middlewares.KeyValue
	err := viper.UnmarshalKey("keys", &keys)
	if err != nil {
		log.Fatalf("Error unmarshaling 'keys', %s", err)
	}

	// Print out the values
	for _, key := range keys {
		if key.Key == "error_code" {
			middlewares.ErrorCodeKey = middlewares.KeyValue{Key: key.Key, Value: key.Value}
		}
	}
}
