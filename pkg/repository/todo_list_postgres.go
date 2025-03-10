package repository

import (
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoGoalPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoGoalPostgres {
	return &TodoGoalPostgres{db: db}
}

func (r *TodoGoalPostgres) Create(userId int, list todo.TodoGoal) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description, end_date) VALUES ($1, $2, $3) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description, list.EndDate)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
func (r *TodoGoalPostgres) GetAll(userId int) ([]todo.TodoGoal, error) {
	var lists []todo.TodoGoal

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoGoalPostgres) GetById(userId, listId int) (todo.TodoGoal, error) {
	var list todo.TodoGoal

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl 
                                       INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}
func (r *TodoGoalPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE ul.list_id = $1 AND ul.user_id = $2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, listId, userId)

	return err
}
func (r *TodoGoalPostgres) Update(userId, listId int, input todo.UpdateGoalInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.EndDate != nil {
		setValues = append(setValues, fmt.Sprintf("end_date=$%d", argId))
		args = append(args, *input.EndDate)
		argId++
	}

	// title = $1
	// description = $1
	// end_date = $1
	// title = $1, description = $2, end_date = $3
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id = $1 AND ul.user_id = $2",
		todoListsTable, setQuery, usersListsTable)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err
}
