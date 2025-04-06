package service

import (
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/llm"
	"github.com/Gayana5/todo-app/pkg/repository"
	"log"
)

type TodoGoalService struct {
	repo repository.TodoGoal
	ai   llm.LLM
}

func NewTodoGoalService(repo repository.TodoGoal, ai llm.LLM) *TodoGoalService {
	return &TodoGoalService{
		repo: repo,
		ai:   ai,
	}
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
func (s *TodoGoalService) AskAI(userId, goalId int) (string, error) {
	goal, err := s.repo.GetById(userId, goalId)
	if err != nil {
		return "", err
	}
	message, err := s.ai.GetAdvice(goal.Title, goal.Description)
	if err != nil {
		return "", err
	}

	return message, nil
}
