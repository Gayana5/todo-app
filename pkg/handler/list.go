package handler

import (
	"database/sql"
	"errors"
	"github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createGoal(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input todo.TodoGoal
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoGoal.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllGoals(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	goals, err := h.services.TodoGoal.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, goals)
}
func (h *Handler) getGoalById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid User Id")
		return
	}
	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Goal Id")
		return
	}
	list, err := h.services.TodoGoal.GetById(userId, goalId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, list)
	}
}
func (h *Handler) updateGoal(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid goal id")
		return
	}

	var input todo.UpdateGoalInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoGoal.Update(userId, goalId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
func (h *Handler) deleteGoal(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid goal id")
		return
	}
	err = h.services.TodoGoal.Delete(userId, goalId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})
}
