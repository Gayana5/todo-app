package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoGoal
}

func NewTodoListService(repo repository.TodoGoal) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoGoal) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoGoal, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoGoal, error) {
	return s.repo.GetById(userId, listId)
}
func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}
func (s *TodoListService) Update(userId, listId int, input todo.UpdateGoalInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
