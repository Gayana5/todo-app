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
}
type TodoGoal interface {
	Create(userId int, goal todo.TodoGoal) (int, error)
	GetAll(userId int) ([]todo.TodoGoal, error)
	GetById(userId, goalId int) (todo.TodoGoal, error)
	Delete(userId, goalId int) error
	Update(userId, goalId int, input todo.UpdateGoalInput) error
}

type TodoItem interface {
	Create(goalId int, item todo.TodoItem) (int, error)
	GetAll(userId, goalId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}
type Repository struct {
	Authorization
	TodoGoal
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoGoal:      NewTodoGoalPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
