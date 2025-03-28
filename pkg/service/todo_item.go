package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	goalRepo repository.TodoGoal
}

func NewTodoItemService(repo repository.TodoItem, goalRepo repository.TodoGoal) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		goalRepo: goalRepo,
	}
}
func (s *TodoItemService) Create(userId, goalId int, item todo.TodoItem) (int, error) {
	if goalId == 0 {
		return s.repo.Create(userId, goalId, item)
	}
	_, err := s.goalRepo.GetById(userId, goalId)
	if err != nil {
		// Список не существует или принадлежит другому пользователю
		return 0, err
	}
	return s.repo.Create(userId, goalId, item)
}
func (s *TodoItemService) GetAll(userId, goalId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, goalId)
}

func (s *TodoItemService) GetById(userId, itemId, goalId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId, goalId)
}

func (s *TodoItemService) Delete(userId, itemId, goalId int) error {
	return s.repo.Delete(userId, itemId, goalId)
}

func (s *TodoItemService) Update(userId, itemId, goalId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, goalId, input)
}
