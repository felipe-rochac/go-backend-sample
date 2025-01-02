package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type DatabaseConfiguration struct {
	Host         string `json,yaml:host`
	Database     string `json,yaml:database`
	User         string `json,yaml:user`
	Password     string `json,yaml:password`
	Port         int    `json,yaml:port`
	MaxLifetime  int    `json,yaml:maxLifetime`
	MaxOpenConns int    `json,yaml:maxOpenConns`
	MaxIdleConns int    `json,yaml:maxIdleConns`
}

var Configuration DatabaseConfiguration

type MySqlDatabase struct {
}

func (db *MySqlDatabase) GetConnection() *sql.DB {
	if reflect.ValueOf(Configuration).IsNil() {
		panic(fmt.Errorf("Database configuration is not initialized."))
	}

	cn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&autocommit=true", Configuration.User, Configuration.Password, Configuration.Host, Configuration.Port, Configuration.Database))

	if err != nil {
		panic(err)
	}

	cn.SetConnMaxLifetime(time.Duration(Configuration.MaxLifetime))
	cn.SetMaxOpenConns(Configuration.MaxOpenConns)
	cn.SetMaxIdleConns(Configuration.MaxIdleConns)

	return cn
}
