package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
)

type TodoGoalService struct {
	repo repository.TodoGoal
}

func NewTodoListService(repo repository.TodoGoal) *TodoGoalService {
	return &TodoGoalService{repo: repo}
}

func (s *TodoGoalService) Create(userId int, goal todo.TodoGoal) (int, error) {
	return s.repo.Create(userId, goal)
}

func (s *TodoGoalService) GetAll(userId int) ([]todo.TodoGoal, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoGoalService) GetById(userId, goalId int) (todo.TodoGoal, error) {
	return s.repo.GetById(userId, goalId)
}
func (s *TodoGoalService) Delete(userId, goalId int) error {
	return s.repo.Delete(userId, goalId)
}
func (s *TodoGoalService) Update(userId, goalId int, input todo.UpdateGoalInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, goalId, input)
}
