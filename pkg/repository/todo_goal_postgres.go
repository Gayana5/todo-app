package repository

import (
	"errors"
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoGoalPostgres struct {
	db *sqlx.DB
}

func NewTodoGoalPostgres(db *sqlx.DB) *TodoGoalPostgres {
	return &TodoGoalPostgres{db: db}
}

func (r *TodoGoalPostgres) Create(userId int, goal todo.TodoGoal) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description, colour) VALUES ($1, $2, $3) RETURNING id", todoGoalsTable)
	row := tx.QueryRow(createListQuery, goal.Title, goal.Description, goal.Colour)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersGoalQuery := fmt.Sprintf("INSERT INTO %s (user_id, goal_id) VALUES ($1, $2)", usersGoalsTable)
	_, err = tx.Exec(createUsersGoalQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
func (r *TodoGoalPostgres) GetAll(userId int) ([]todo.TodoGoal, error) {
	var goals []todo.TodoGoal

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, tl.colour, tl.completed_tasks, tl.total_tasks FROM %s tl INNER JOIN %s ul on tl.id = ul.goal_id WHERE ul.user_id = $1",
		todoGoalsTable, usersGoalsTable)
	err := r.db.Select(&goals, query, userId)
	if err != nil {
		return goals, err
	}
	for i := range goals {
		goal := &goals[i]
		if goal.TotalTasks != 0 {
			goal.Progress = int((float64(goal.CompletedTasks) / float64(goal.TotalTasks)) * 100)
		}
	}

	return goals, nil
}

func (r *TodoGoalPostgres) GetById(userId, goalId int) (todo.TodoGoal, error) {
	var goal todo.TodoGoal

	query := fmt.Sprintf(`
        SELECT tl.id, tl.title, tl.description, tl.colour, tl.completed_tasks, tl.total_tasks
        FROM %s tl
        INNER JOIN %s ul ON tl.id = ul.goal_id
        WHERE ul.user_id = $1 AND ul.goal_id = $2
    `, todoGoalsTable, usersGoalsTable)
	err := r.db.Get(&goal, query, userId, goalId)
	if err != nil {
		return goal, err
	}

	if goal.TotalTasks != 0 {
		goal.Progress = int((float64(goal.CompletedTasks) / float64(goal.TotalTasks)) * 100)

	}
	return goal, nil
}
func (r *TodoGoalPostgres) Delete(userId, goalId int) error {
	deleteTaskquery := fmt.Sprintf(`DELETE FROM todo_tasks WHERE goal_id = $1`)
	_, err := r.db.Exec(deleteTaskquery, goalId)
	if err != nil {
		return err
	}

	deleteGoalQuery := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE id = $1 
        AND id IN (SELECT goal_id FROM %s WHERE user_id = $2)
    `, todoGoalsTable, usersGoalsTable)

	_, err = r.db.Exec(deleteGoalQuery, goalId, userId)
	return err
}
func (r *TodoGoalPostgres) Update(userId, goalId int, input todo.UpdateGoalInput) error {
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
	if input.Colour != nil {
		setValues = append(setValues, fmt.Sprintf("colour=$%d", argId))
		args = append(args, *input.Colour)
		argId++
	}

	if len(setValues) == 0 {
		return errors.New("update structure has no values")
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
    UPDATE %s tl 
    SET %s 
    WHERE tl.id IN (
        SELECT ul.goal_id FROM %s ul WHERE ul.goal_id = $%d AND ul.user_id = $%d
    )`,
		todoGoalsTable, setQuery, usersGoalsTable, argId, argId+1,
	)

	args = append(args, goalId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err
}
