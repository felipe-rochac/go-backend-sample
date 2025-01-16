package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type MockDatabase struct {
}

var dbConfig = DatabaseConfiguration{
	Host:         "host",
	Database:     "database",
	User:         "user",
	Password:     "password",
	Port:         3306,
	MaxLifetime:  10,
	MaxOpenConns: 100,
	MaxIdleConns: 10,
}

var cn *sql.DB
var mock sqlmock.Sqlmock

func TestMain_Config(m *testing.M) {
	var err error
	cn, mock, err = sqlmock.New()

	if err != nil {
		panic(err)
	}

	// Setup code before running tests
	code := m.Run()

	// Teardown code after running tests
	os.Exit(code)
}
func Test_GetConnection_ExpectSucces(t *testing.T) {
	tests := []struct {
		name           string
		config         DatabaseConfiguration
		expectedErrMsg string
	}{
		{
			name:           "Empty configuration",
			config:         DatabaseConfiguration{},
			expectedErrMsg: "database configuration is not initialized.",
		},
		{
			name: "Invalid connection string",
			config: DatabaseConfiguration{
				Host:         "invalid_host",
				Database:     "database",
				User:         "user",
				Password:     "password",
				Port:         3306,
				MaxLifetime:  10,
				MaxOpenConns: 100,
				MaxIdleConns: 10,
			},
			expectedErrMsg: "could not open connection to host invalid_host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &MySqlDatabaseService{
				Configuration: tt.config,
			}

			_, err := db.GetConnection()
			if err == nil {
				t.Fatalf("expected an error but got nil")
			}

			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func Test_GetConnection_ExpectSuccess(t *testing.T) {
	db := &MySqlDatabaseService{
		Configuration: dbConfig,
	}

	conn, err := db.GetConnection()
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}
