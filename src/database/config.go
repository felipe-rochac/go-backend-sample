package database

import (
	"backend-sample/common"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

type MySqlDatabaseService struct {
	Configuration DatabaseConfiguration
	db            *sql.DB
}

func (m *MySqlDatabaseService) GetConnection() (*sql.DB, *common.BackendError) {
	isEmpty := m.Configuration == DatabaseConfiguration{}
	if isEmpty {
		return nil, common.NewBackendError(500, "GetConnection.1", "database configuration is not initialized.", nil)
	}

	if m.db == nil {
		var err error
		m.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&autocommit=true", m.Configuration.User, m.Configuration.Password, m.Configuration.Host, m.Configuration.Port, m.Configuration.Database))
		if err != nil {
			return nil, common.NewBackendError(500, "GetConnection.2", "could not open connection to host %s", err, m.Configuration.Host)
		}
	}

	m.db.SetConnMaxLifetime(time.Duration(m.Configuration.MaxLifetime))
	m.db.SetMaxOpenConns(m.Configuration.MaxOpenConns)
	m.db.SetMaxIdleConns(m.Configuration.MaxIdleConns)

	return m.db, nil
}
