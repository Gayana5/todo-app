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
	Create(userId int, goal todo.TodoGoal) (int, error)
	GetAll(userId int) ([]todo.TodoGoal, error)
	GetById(userId, goalId int) (todo.TodoGoal, error)
	Delete(userId, goalId int) error
	Update(userId, goalId int, input todo.UpdateGoalInput) error
}

type TodoItem interface {
	Create(userId, goalId int, item todo.TodoItem) (int, error)
	GetAll(userId, goalId int) ([]todo.TodoItem, error)
	GetById(userId, itemId, goalId int) (todo.TodoItem, error)
	Delete(userId, itemId, goalId int) error
	Update(userId, itemId, goalId int, input todo.UpdateItemInput) error
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
