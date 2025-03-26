package repository

import (
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
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
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	if err != nil {
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
