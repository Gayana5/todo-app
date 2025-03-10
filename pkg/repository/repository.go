package repository

import (
	"github.com/Gayana5/todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
	UserExists(email string) (bool, error)
}
type TodoGoal interface {
	Create(userId int, list todo.TodoGoal) (int, error)
	GetAll(userId int) ([]todo.TodoGoal, error)
	GetById(userId, listId int) (todo.TodoGoal, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateGoalInput) error
}

type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
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
		TodoGoal:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
