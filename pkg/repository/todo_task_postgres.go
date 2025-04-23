package repository

import (
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoTaskPostgres struct {
	db *sqlx.DB
}

func NewTodoTaskPostgres(db *sqlx.DB) *TodoTaskPostgres {
	return &TodoTaskPostgres{db: db}
}
func (r *TodoTaskPostgres) Create(userId int, goalId int, task todo.TodoTask) (int, error) {
	if r.db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var taskId int
	createTaskQuery := fmt.Sprintf(
		`INSERT INTO %s (user_id, title, description, goal_id, end_date, start_time, end_time, colour, done) 
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		todoTasksTable,
	)
	row := tx.QueryRow(createTaskQuery, userId, task.Title, task.Description, goalId, task.EndDate, task.StartTime, task.EndTime, task.Colour, task.Done)
	err = row.Scan(&taskId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	if goalId != 0 {
		// Добавляем связь задачи с целью
		createGoalTasksQuery := fmt.Sprintf("INSERT INTO %s (goal_id, task_id) VALUES ($1, $2)", goalsTaskTable)
		_, err = tx.Exec(createGoalTasksQuery, goalId, taskId)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return 0, err
			}
			return 0, err
		}

		// Обновляем total_tasks в цели
		updateGoalQuery := fmt.Sprintf(
			"UPDATE %s SET total_tasks = total_tasks + 1 WHERE id = $1",
			todoGoalsTable,
		)
		_, err = tx.Exec(updateGoalQuery, goalId)
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
	return taskId, nil
}

func (r *TodoTaskPostgres) GetAll(userId, goalId int) ([]todo.TodoTask, error) {
	var tasks []todo.TodoTask
	if goalId != 0 {
		query := fmt.Sprintf(
			`SELECT ti.id, ti.user_id, ti.title, ti.description, ti.goal_id, ti.end_date, ti.start_time, ti.end_time, ti.colour, ti.done 
				FROM %s ti INNER JOIN %s li ON li.task_id = ti.id 
				INNER JOIN %s ul ON ul.goal_id = li.goal_id 
				WHERE li.goal_id = $1 AND ul.user_id = $2`,
			todoTasksTable, goalsTaskTable, usersGoalsTable,
		)

		if err := r.db.Select(&tasks, query, goalId, userId); err != nil {
			return tasks, err
		}
	} else {
		query := fmt.Sprintf(
			`SELECT id, user_id, title, description, goal_id, end_date, start_time, end_time, colour, done 
     				FROM %s 
     				WHERE user_id = $1`,
			todoTasksTable,
		)

		if err := r.db.Select(&tasks, query, userId); err != nil {
			return tasks, err
		}
	}
	return tasks, nil
}

func (r *TodoTaskPostgres) GetById(userId, taskId, goalId int) (todo.TodoTask, error) {
	var task todo.TodoTask
	if goalId != 0 {
		query := fmt.Sprintf(`SELECT ti.id, ti.user_id, ti.title, ti.description, ti.goal_id, ti.end_date, ti.start_time, ti.end_time, ti.colour, ti.done
									FROM %s ti INNER JOIN %s li on li.task_id = ti.id
									INNER JOIN %s ul on ul.goal_id = li.goal_id WHERE ti.id = $1 AND ul.user_id = $2`,
			todoTasksTable, goalsTaskTable, usersGoalsTable)
		if err := r.db.Get(&task, query, taskId, userId); err != nil {
			return task, err
		}
	} else {
		query := fmt.Sprintf(`
		    SELECT id, user_id, title, description, goal_id, end_date, start_time, end_time, colour, done
		    FROM %s
		    WHERE id = $1 AND user_id = $2`,
			todoTasksTable)

		if err := r.db.Get(&task, query, taskId, userId); err != nil {
			return task, err
		}
	}

	return task, nil
}

func (r *TodoTaskPostgres) Delete(userId, taskId, goalId int) error {
	if goalId != 0 {
		query := fmt.Sprintf(
			`DELETE FROM %s ti 
        USING %s li, %s ul 
        WHERE ti.id = li.task_id 
        AND li.goal_id = ul.goal_id 
        AND ul.user_id = $1 
        AND ti.id = $2`,
			todoTasksTable, goalsTaskTable, usersGoalsTable,
		)

		_, err := r.db.Exec(query, userId, taskId)
		return err
	} else {
		query := fmt.Sprintf(
			`DELETE FROM %s ti  
                    WHERE ti.user_id = $1 AND ti.id = $2`,
			todoTasksTable,
		)

		_, err := r.db.Exec(query, userId, taskId)
		return err
	}
}

func (r *TodoTaskPostgres) Update(userId, taskId, goalId int, input todo.UpdateTaskInput) error {
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
	if input.Colour != nil {
		setValues = append(setValues, fmt.Sprintf("colour=$%d", argId))
		args = append(args, *input.Colour)
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

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Получаем текущий статус и цель задачи
	var currentDone bool
	var currentGoalId int
	err = tx.QueryRow(
		fmt.Sprintf("SELECT done, goal_id FROM %s WHERE id = $1 AND user_id = $2", todoTasksTable),
		taskId, userId,
	).Scan(&currentDone, &currentGoalId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	// Выполняем обновление задачи
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		`UPDATE %s ti SET %s WHERE ti.id = $%d AND ti.user_id = $%d`,
		todoTasksTable, setQuery, argId, argId+1,
	)
	args = append(args, taskId, userId)

	_, err = tx.Exec(query, args...)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if input.Done != nil {
		newDone := *input.Done
		if newDone != currentDone && currentGoalId != 0 {
			delta := 1
			if !newDone {
				delta = -1
			}

			updateGoalQuery := fmt.Sprintf(
				"UPDATE %s SET completed_tasks = completed_tasks + $1 WHERE id = $2",
				todoGoalsTable,
			)
			_, err = tx.Exec(updateGoalQuery, delta, currentGoalId)
			if err != nil {
				err := tx.Rollback()
				if err != nil {
					return err
				}
				return err
			}
		}
	}

	return tx.Commit()
}
