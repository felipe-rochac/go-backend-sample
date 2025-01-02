package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	Id                    uuid.UUID
	Name, Email, Password string
}

type UserWhereClause struct {
	Ids         []uuid.UUID
	Name, Email string
}

var db MySqlDatabase
var insertUserQuery string = `INSERT INTO user (user_id, name, email, password) VALUES (?, ?, ?, ?)`
var updateUserQuery string = `UPDATE user SET name = ?, email = ?, password = ? WHERE user_id = ?`
var deleteUserQuery string = `DELETE FROM users WHERE id = ?`

func createUser(user *User) (uuid.UUID, bool) {
	cn := db.GetConnection()
	defer cn.Close()
	user.Id = uuid.New()

	_, err := cn.Exec(insertUserQuery, user.Id, user.Name, user.Email, user.Email, user.Password)

	if err != nil {
		log.Fatal(err)
		return user.Id, false
	}

	return user.Id, true
}

func updateUser(user User) bool {
	cn := db.GetConnection()
	defer cn.Close()
	user.Id = uuid.New()

	_, err := cn.Exec(updateUserQuery, user.Id, user.Name, user.Email, user.Email, user.Password)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func GetUsers(where UserWhereClause) ([]User, bool) {
	users := make([]User, 0)
	cn := db.GetConnection()
	defer cn.Close()
	query := "SELECT user_id, name, email, password FROM users "
	clause, values := buildWhereClause(where)
	if len(clause) > 0 {
		query += " WHERE " + clause
	}

	rows, err := cn.Query(query, values...)

	if err != nil {
		log.Fatal(err)
		return nil, false
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
			return []User{}, false
		}

		users = append(users, User{Id: uuid, Name: name, Email: email, Password: password})
	}

	return users, true
}

func DeleteUser(uuid uuid.UUIDs) bool {
	cn := db.GetConnection()
	defer cn.Close()

	_, err := cn.Exec(deleteUserQuery, uuid)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
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
