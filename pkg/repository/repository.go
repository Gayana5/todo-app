package repository

import (
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
	UserExists(email string) (bool, error)
	GetInfo(id int) (todo.User, error)
	UpdateInfo(userId int, input todo.UpdateUserInput) error
	ResetPassword(email, password string) error
}
type TodoGoal interface {
	Create(userId int, goal todo.TodoGoal) (int, error)
	GetAll(userId int) ([]todo.TodoGoal, error)
	GetById(userId, goalId int) (todo.TodoGoal, error)
	Delete(userId, goalId int) error
	Update(userId, goalId int, input todo.UpdateGoalInput) error
}

type TodoTask interface {
	Create(userId int, goalId int, task todo.TodoTask) (int, error)
	GetAll(userId, goalId int) ([]todo.TodoTask, error)
	GetById(userId, taskId, goalId int) (todo.TodoTask, error)
	Delete(userId, taskId, goalId int) error
	Update(userId, taskId, goalId int, input todo.UpdateTaskInput) error
}
type Repository struct {
	Authorization
	TodoGoal
	TodoTask
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoGoal:      NewTodoGoalPostgres(db),
		TodoTask:      NewTodoTaskPostgres(db),
	}
}
