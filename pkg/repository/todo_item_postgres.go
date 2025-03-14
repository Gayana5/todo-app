package repository

import (
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}
func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {

	if r.db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description, end_date, start_time, end_time, priority, done) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		todoItemsTable,
	)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description, item.EndDate, item.StartTime, item.EndTime, item.Priority, item.Done)

	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemTable)

	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return itemId, nil
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.end_date, ti.start_time, ti.end_time, ti.priority, ti.done FROM %s ti "+
			"INNER JOIN %s li ON li.item_id = ti.id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE li.list_id = $1 AND ul.user_id = $2",
		todoItemsTable, listsItemTable, usersListsTable,
	)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.end_date, ti.start_time, ti.end_time, ti.priority, ti.done FROM %s ti "+
			"INNER JOIN %s li ON li.item_id = ti.id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE ti.id = $1 AND ul.user_id = $2",
		todoItemsTable, listsItemTable, usersListsTable,
	)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s ti 
        USING %s li, %s ul 
        WHERE ti.id = li.item_id 
        AND li.list_id = ul.list_id 
        AND ul.user_id = $1 
        AND ti.id = $2`,
		todoItemsTable, listsItemTable, usersListsTable,
	)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
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
	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}
	if input.StartTime != nil {
		setValues = append(setValues, fmt.Sprintf("start_time=$%d", argId))
		args = append(args, *input.StartTime)
		argId++
	}
	if input.EndTime != nil {
		setValues = append(setValues, fmt.Sprintf("end_time=$%d", argId))
		args = append(args, *input.EndTime)
		argId++
	}
	if input.Priority != nil {
		setValues = append(setValues, fmt.Sprintf("priority=$%d", argId))
		args = append(args, *input.Priority)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		"UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d",
		todoItemsTable, setQuery, listsItemTable, usersListsTable, argId, argId+1,
	)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
