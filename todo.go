package todo

import (
	"errors"
	"time"
)

type TodoGoal struct {
	Id             int    `json:"id" db:"id"`
	Title          string `json:"title" db:"title" binding:"required"`
	Description    string `json:"description" db:"description"`
	Colour         int    `json:"colour" db:"colour" binding:"required"`
	Progress       int    `json:"progress" db:"progress"`
	CompletedTasks int    `json:"completed_tasks" db:"completed_tasks"`
	TotalTasks     int    `json:"total_tasks" db:"total_tasks"`
}

type TodoTask struct {
	Id          int       `json:"id" db:"id"`
	UserId      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	GoalId      int       `json:"goal_id" db:"goal_id"`
	EndDate     time.Time `json:"end_date" db:"end_date" binding:"required"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Colour      int       `json:"colour" db:"colour"`
	Done        bool      `json:"done" db:"done"`
}

type UpdateGoalInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Colour      *int    `json:"colour"`
}

func (i UpdateGoalInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Colour == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdateTaskInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	EndDate     *time.Time `json:"end_date"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Colour      *int       `json:"colour"`
	Done        *bool      `json:"done"`
}

func (i UpdateTaskInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil && i.EndDate == nil && i.StartTime == nil && i.EndTime == nil && i.Colour == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
