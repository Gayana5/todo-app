//go:generate mockgen -source=service.go -destination=mocks/mock.go

package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/llm"
	"github.com/Gayana5/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GenerateCode() string
	GetInfo(id int) (todo.User, error)
	UpdateInfo(userId int, input todo.UpdateUserInput) error
	UserExists(email string) (bool, error)
	ResetPassword(email, password string) error
	SendCodeToEmail(to string, code string) error
}
type TodoGoal interface {
	Create(userId int, goal todo.TodoGoal) (int, error)
	GetAll(userId int) ([]todo.TodoGoal, error)
	GetById(userId, goalId int) (todo.TodoGoal, error)
	Delete(userId, goalId int) error
	Update(userId, goalId int, input todo.UpdateGoalInput) error
	AskAI(userId, goalId int) (string, error)
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

func NewService(repos *repository.Repository, ai llm.LLM) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoGoal:      NewTodoGoalService(repos.TodoGoal, ai),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoGoal),
	}
}
