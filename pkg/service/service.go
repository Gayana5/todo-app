package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GenerateCode() string
	GetInfo(id int) (todo.User, error)
}
type TodoGoal interface {
	Create(userId int, list todo.TodoGoal) (int, error)
	GetAll(userId int) ([]todo.TodoGoal, error)
	GetById(userId, listId int) (todo.TodoGoal, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateGoalInput) error
}

type TodoItem interface {
	Create(userId, listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}
type Service struct {
	Authorization
	TodoGoal
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoGoal:      NewTodoListService(repos.TodoGoal),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoGoal),
	}
}
