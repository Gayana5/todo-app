package handler

import (
	"database/sql"
	"errors"
	"github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createTask(c *gin.Context) {
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
	var input todo.TodoTask
	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoTask.Create(userId, goalId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}
func (h *Handler) getAllTasks(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid User Id")
		return
	}
	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	tasks, err := h.services.TodoTask.GetAll(userId, goalId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
}
func (h *Handler) getTaskById(c *gin.Context) {
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
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Task Id")
		return
	}

	task, err := h.services.TodoTask.GetById(userId, taskId, goalId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, task)
	}
}
func (h *Handler) updateTask(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Task id")
		return
	}
	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Goal id")
		return
	}

	var input todo.UpdateTaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoTask.Update(userId, taskId, goalId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
func (h *Handler) deleteTask(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid User id")
		return
	}
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Task id")
		return
	}
	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Goal id")
		return
	}
	err = h.services.TodoTask.Delete(userId, taskId, goalId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}
