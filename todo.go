package todo

import (
	"errors"
	"time"
)

type TodoGoal struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Colour      int    `json:"colour" db:"colour" binding:"required"`
	Progress    int    `json:"progress" db:"progress"`
}

type UsersGoals struct {
	Id     int
	UserId int
	ListId int
}
type TodoItem struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	EndDate     time.Time `json:"end_date" db:"end_date" binding:"required"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Priority    bool      `json:"priority" db:"priority" binding:"required"`
	Done        bool      `json:"done" db:"done"`
}
type GoalsItem struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	EndDate     time.Time `json:"end_date" db:"end_date" binding:"required"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Priority    bool      `json:"priority" db:"priority" binding:"required"`
	Done        bool      `json:"done" db:"done"`
}

type UpdateGoalInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Colour      *int    `json:"colour"`
}

func (i UpdateGoalInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdateItemInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date" db:"date"`
	StartTime   *time.Time `json:"start_time" db:"start_time"`
	EndTime     *time.Time `json:"end_time" db:"end_time"`
	Priority    *bool      `json:"priority" db:"priority"`
	Done        *bool      `json:"done" db:"done"`
}

func (i UpdateItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
