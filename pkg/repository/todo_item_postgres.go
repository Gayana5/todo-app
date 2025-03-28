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
func (r *TodoItemPostgres) Create(userId int, goalId int, item todo.TodoItem) (int, error) {

	if r.db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf(
		`INSERT INTO %s (user_id, title, description, end_date, start_time, end_time, priority, done) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		todoItemsTable,
	)

	row := tx.QueryRow(createItemQuery, userId, item.Title, item.Description, item.EndDate, item.StartTime, item.EndTime, item.Priority, item.Done)

	err = row.Scan(&itemId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}
	if goalId != 0 {
		createGoalItemsQuery := fmt.Sprintf("INSERT INTO %s (goal_id, item_id) VALUES ($1, $2)", goalsItemTable)

		_, err = tx.Exec(createGoalItemsQuery, goalId, itemId)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return itemId, nil
}

func (r *TodoItemPostgres) GetAll(userId, goalId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	if goalId != 0 {
		query := fmt.Sprintf(
			`SELECT ti.id, ti.user_id, ti.title, ti.description, ti.end_date, ti.start_time, ti.end_time, ti.priority, ti.done 
				FROM %s ti INNER JOIN %s li ON li.item_id = ti.id 
				INNER JOIN %s ul ON ul.goal_id = li.goal_id 
				WHERE li.goal_id = $1 AND ul.user_id = $2`,
			todoItemsTable, goalsItemTable, usersGoalsTable,
		)

		if err := r.db.Select(&items, query, goalId, userId); err != nil {
			return nil, err
		}
	} else {
		query := fmt.Sprintf(
			`SELECT id, user_id, title, description, end_date, start_time, end_time, priority, done 
     				FROM %s 
     				WHERE user_id = $1`,
			todoItemsTable,
		)

		if err := r.db.Select(&items, query, userId); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId, goalId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	if goalId != 0 {
		query := fmt.Sprintf(`SELECT ti.id, ti.user_id, ti.title, ti.description, ti.end_date, ti.start_time, ti.end_time, ti.priority, ti.done
									FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.goal_id = li.goal_id WHERE ti.id = $1 AND ul.user_id = $2`,
			todoItemsTable, goalsItemTable, usersGoalsTable)
		if err := r.db.Get(&item, query, itemId, userId); err != nil {
			return item, err
		}
	} else {
		query := fmt.Sprintf(`
		    SELECT id, user_id, title, description, end_date, start_time, end_time, priority, done
		    FROM %s
		    WHERE id = $1 AND user_id = $2`,
			todoItemsTable)

		if err := r.db.Get(&item, query, itemId, userId); err != nil {
			return item, err
		}
	}

	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s ti 
        USING %s li, %s ul 
        WHERE ti.id = li.item_id 
        AND li.goal_id = ul.goal_id 
        AND ul.user_id = $1 
        AND ti.id = $2`,
		todoItemsTable, goalsItemTable, usersGoalsTable,
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
		`UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.goal_id = ul.goal_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, goalsItemTable, usersGoalsTable, argId, argId+1,
	)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
