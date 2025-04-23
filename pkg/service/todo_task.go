package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
)

type TodoTaskService struct {
	repo     repository.TodoTask
	goalRepo repository.TodoGoal
}

func NewTodoTaskService(repo repository.TodoTask, goalRepo repository.TodoGoal) *TodoTaskService {
	return &TodoTaskService{
		repo:     repo,
		goalRepo: goalRepo,
	}
}
func (s *TodoTaskService) Create(userId, goalId int, task todo.TodoTask) (int, error) {
	if goalId == 0 {
		return s.repo.Create(userId, goalId, task)
	}
	_, err := s.goalRepo.GetById(userId, goalId)
	if err != nil {
		// Список не существует или принадлежит другому пользователю
		return 0, err
	}
	return s.repo.Create(userId, goalId, task)
}
func (s *TodoTaskService) GetAll(userId, goalId int) ([]todo.TodoTask, error) {
	return s.repo.GetAll(userId, goalId)
}

func (s *TodoTaskService) GetById(userId, taskId, goalId int) (todo.TodoTask, error) {
	return s.repo.GetById(userId, taskId, goalId)
}

func (s *TodoTaskService) Delete(userId, taskId, goalId int) error {
	return s.repo.Delete(userId, taskId, goalId)
}

func (s *TodoTaskService) Update(userId, taskId, goalId int, input todo.UpdateTaskInput) error {
	return s.repo.Update(userId, taskId, goalId, input)
}
