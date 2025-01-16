package database

import (
	"backend-sample/common"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type UserEntity struct {
	Id                    uuid.UUID
	Name, Email, Password string
}

type UserWhereClause struct {
	Ids         []uuid.UUID
	Name, Email string
}

var (
	insertUserQuery     string = `INSERT INTO user (user_id, name, email, password) VALUES (?, ?, ?, ?)`
	updateUserQuery     string = `UPDATE user SET name = ?, email = ?, password = ? WHERE user_id = ?`
	deleteUserQuery     string = `DELETE FROM user WHERE user_id = ?`
	selectUserByIdQuery string = `SELECT user_id, name, email, pasword FROM user WHERE user_id = ?`
)

type UsersRepository interface {
	CreateUser(user *UserEntity) (*UserEntity, *common.BackendError)
	UpdateUser(user UserEntity) *common.BackendError
	GetUsers(where UserWhereClause) (*[]UserEntity, *common.BackendError)
	GetUsersByName(name string, exactMatch bool) (*[]UserEntity, *common.BackendError)
	GetUserById(uuid uuid.UUID) (*UserEntity, *common.BackendError)
	DeleteUser(uuid uuid.UUID) *common.BackendError
}

type repositoryService struct {
	db MySqlDatabaseService
}

func (repo *repositoryService) CreateUser(name, email, password string) (*UserEntity, *common.BackendError) {
	cn, berr := repo.db.GetConnection()

	if berr != nil {
		return nil, berr
	}
	defer cn.Close()

	id := uuid.New()
	binary, err := common.UuidToBinary(id)

	if err != nil {
		return nil, common.NewBackendError(500, "CreateUser.1", "could generate an uuid", err)
	}

	_, err = cn.Exec(insertUserQuery, binary, name, email, password)

	if err != nil {
		return nil, common.NewBackendError(500, "CreateUser.2", "could not insert user", err)
	}

	user, err := repo.GetUserById(id)

	if err != nil {
		return nil, common.NewBackendError(500, "CreateUser.3", "could not retrieve user %s", err, name)
	}

	return user, nil
}

func (repo *repositoryService) UpdateUser(user UserEntity) *common.BackendError {
	cn, berr := repo.db.GetConnection()

	if berr != nil {
		return berr
	}
	defer cn.Close()

	id, err := common.UuidToBinary(user.Id)

	if err != nil {
		return common.NewBackendError(500, "UpdateUser.1", "error converting uuid to binary.", err)
	}

	result, err := cn.Exec(updateUserQuery, user.Name, user.Email, user.Password, id)

	if err != nil {
		return common.NewBackendError(500, "UpdateUser.2", "error executing query.", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return common.NewBackendError(500, "UpdateUser.3", "error reading rows.", err)
	}

	if rowsAffected == 0 {
		log.Println("No rows found")
		return nil
	}

	return nil
}

func (repo *repositoryService) GetUsersByName(name string, exactMatch bool) (*[]UserEntity, *common.BackendError) {
	cn, berr := repo.db.GetConnection()

	if berr != nil {
		return nil, berr
	}

	defer cn.Close()

	var operator string = "="
	if !exactMatch {
		operator = "like"
	}
	rows, err := cn.Query(fmt.Sprintf("SELECT user_id, name, email, pasword FROM user WHERE name %s ?", operator), name)

	if err != nil {
		return nil, common.NewBackendError(500, "GetUserByName.1", "error querying user by name %s.", err, name)
	}

	users := make([]UserEntity, 0)
	for rows.Next() {
		var id []byte
		var email, password string
		err = rows.Scan(&id, &name, &email, &password)
		if err != nil {
			return nil, common.NewBackendError(500, "GetUserByName.2", "error reading row.", err, name)
		}

		uuid, err := uuid.FromBytes(id)

		if err != nil {
			return nil, common.NewBackendError(500, "GetUserByName.3", "error parsing user id to uuid.", err)
		}

		users = append(users, UserEntity{Id: uuid, Name: name, Email: email, Password: password})
	}

	return &users, nil

}

func (repo *repositoryService) GetUserById(id uuid.UUID) (*UserEntity, *common.BackendError) {
	cn, berr := repo.db.GetConnection()

	if berr != nil {
		return nil, berr
	}

	defer cn.Close()

	binary, err := common.UuidToBinary(id)

	if err != nil {
		return nil, common.NewBackendError(500, "GetUserById.1", "converting uuid %s.", err, id.String())
	}

	rows, err := cn.Query(selectUserByIdQuery, binary)

	if err != nil {
		return nil, common.NewBackendError(500, "GetUserById.2", "error querying user by id %s.", err, id.String())
	}

	if !rows.Next() {
		return nil, nil
	}

	var name, email, password string
	err = rows.Scan(&binary, &name, &email, &password)
	if err != nil {
		return nil, common.NewBackendError(500, "GetUserByName.3", "error reading row.", err, name)
	}

	uuid, err := uuid.FromBytes(binary)

	if err != nil {
		return nil, common.NewBackendError(500, "GetUserByName.4", "error parsing user id to uuid.", err)
	}

	return &UserEntity{Id: uuid, Name: name, Email: email, Password: password}, nil

}

func (repo *repositoryService) GetUsers(where UserWhereClause) (*[]UserEntity, *common.BackendError) {
	users := make([]UserEntity, 0)
	cn, berr := repo.db.GetConnection()

	if berr != nil {
		return nil, berr
	}

	defer cn.Close()
	query := "SELECT user_id, name, email, password FROM user "
	clause, values := buildWhereClause(where)
	if len(clause) > 0 {
		query += " WHERE " + clause
	}

	rows, err := cn.Query(query, values...)

	if err != nil {
		return &[]UserEntity{}, common.NewBackendError(500, "GetUsers.1", "could not execute query.", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id []byte
		var name, email, password string
		err := rows.Scan(&id, &name, &email, &password)
		if err != nil {
			log.Fatal(err)
		}

		uuid, err := uuid.FromBytes(id)

		if err != nil {
			log.Fatal(err)
			return &[]UserEntity{}, common.NewBackendError(500, "GetUsers.2", "could not parse id to uuid.", err)
		}

		users = append(users, UserEntity{Id: uuid, Name: name, Email: email, Password: password})
	}

	return &users, nil
}

func (repo *repositoryService) DeleteUser(uuid uuid.UUID) *common.BackendError {
	cn, berr := repo.db.GetConnection()

	if berr != nil {
		return berr
	}

	defer cn.Close()

	id, err := common.UuidToBinary(uuid)

	if err != nil {
		return common.NewBackendError(500, "DeleteUser.1", "cannot parse id to uuid", err)
	}

	result, err := cn.Exec(deleteUserQuery, id)
	if err != nil {
		return common.NewBackendError(500, "DeleteUser.2", "cannot execute query", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return common.NewBackendError(500, "DeleteUser.3", "failed to retrieve affected rows", err)
	}

	if rowsAffected == 0 {
		log.Println("No rows deleted. UUID may not exist.")
	}

	return nil
}

func buildWhereClause(where UserWhereClause) (string, []interface{}) {
	var builder strings.Builder
	var values []interface{}
	conditions := 0

	addCondition := func(clause string, value interface{}) {
		if conditions > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString(clause)
		if value != nil {
			values = append(values, value)
		}
		conditions++
	}

	// Add conditions for IDs
	if len(where.Ids) > 0 {
		placeholders := strings.Repeat("?, ", len(where.Ids))
		placeholders = placeholders[:len(placeholders)-2] // Remove trailing ", "
		addCondition(fmt.Sprintf("user_id IN (%s)", placeholders), nil)
		values = append(values, where.Ids)
	}

	// Add condition for Name
	if len(where.Name) > 0 {
		addCondition("name LIKE ?", where.Name)
	}

	// Add condition for Email
	if len(where.Email) > 0 {
		addCondition("email LIKE ?", where.Email)
	}

	return builder.String(), values
}
