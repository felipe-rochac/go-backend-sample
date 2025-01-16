package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

var repo repositoryService
var sqlCnMock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	var err error
	var db *sql.DB
	db, sqlCnMock, err = sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
	defer db.Close()
	mdb := MySqlDatabaseService{dbConfig, db}
	repo = repositoryService{mdb}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	result := m.Run()

	os.Exit(result)
}

func Test_CreateUser_ExpectSuccess(t *testing.T) {
	sqlCnMock.ExpectExec("INSERT INTO user").
		WithArgs(sqlmock.AnyArg(), "John Doe", "john@example.com", "password").
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlCnMock.ExpectQuery("SELECT user_id, name, email, pasword FROM user WHERE user_id = ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "email", "password"}).
			AddRow([]byte{1, 2, 3, 4}, "John Doe", "john@example.com", "password"))

	user, err := repo.CreateUser("John Doe", "john@example.com", "password")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected user name to be 'John Doe', got '%s'", user.Name)
	}
}

func Test_UpdateUser_ExpectSuccess(t *testing.T) {
	sqlCnMock.ExpectExec("UPDATE user SET name = ?, email = ?, password = ? WHERE user_id = ?").
		WithArgs("John Doe", "john@example.com", "newpassword", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateUser(UserEntity{Id: uuid.New(), Name: "John Doe", Email: "john@example.com", Password: "newpassword"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func Test_GetUsersByName_ExpectSuccess(t *testing.T) {
	sqlCnMock.ExpectQuery("SELECT user_id, name, email, pasword FROM user WHERE name = ?").
		WithArgs("John Doe").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "email", "password"}).
			AddRow([]byte{1, 2, 3, 4}, "John Doe", "john@example.com", "password"))

	users, err := repo.GetUsersByName("John Doe", true)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(*users) != 1 {
		t.Errorf("expected 1 user, got %d", len(*users))
	}
}

func Test_GetUserById_ExpectSuccess(t *testing.T) {
	sqlCnMock.ExpectQuery("SELECT user_id, name, email, pasword FROM user WHERE user_id = ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "email", "password"}).
			AddRow([]byte{1, 2, 3, 4}, "John Doe", "john@example.com", "password"))

	user, err := repo.GetUserById(uuid.New())
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected user name to be 'John Doe', got '%s'", user.Name)
	}
}

func Test_DeleteUser_ExpectSuccess(t *testing.T) {
	sqlCnMock.ExpectExec("DELETE FROM user WHERE user_id = ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteUser(uuid.New())
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}


