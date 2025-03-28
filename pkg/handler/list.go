package handler

import (
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

type getAllGoalsResponse struct {
	Data []todo.TodoGoal `json:"data"`
}

func (h *Handler) getAllGoals(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	goals, err := h.services.TodoGoal.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllGoalsResponse{
		Data: goals,
	})
}
func (h *Handler) getGoalById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	list, err := h.services.TodoGoal.GetById(userId, goalId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
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
