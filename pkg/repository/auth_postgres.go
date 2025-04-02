package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (first_name, second_name, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.FirstName, user.SecondName, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}
func (r *AuthPostgres) UserExists(email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}
func (r *AuthPostgres) GetInfo(id int) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}
func (r *AuthPostgres) UpdateInfo(userId int, input todo.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argId))
		args = append(args, *input.FirstName)
		argId++
	}
	if input.SecondName != nil {
		setValues = append(setValues, fmt.Sprintf("second_name=$%d", argId))
		args = append(args, *input.SecondName)
		argId++
	}
	if len(setValues) == 0 {
		return errors.New("update structure has no values")
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=$%d",
		usersTable, setQuery, argId,
	)

	args = append(args, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err
}
func (r *AuthPostgres) ResetPassword(email, password string) error {
	if password == "" {
		return errors.New("invalid password")
	}
	query := fmt.Sprintf(
		"UPDATE %s SET password_hash = $1 WHERE email = $2",
		usersTable,
	)
	_, err := r.db.Exec(query, password, email)
	return err
}
