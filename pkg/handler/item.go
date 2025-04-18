package handler

import (
	"database/sql"
	"errors"
	"github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createItem(c *gin.Context) {
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
	var input todo.TodoItem
	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoItem.Create(userId, goalId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}
func (h *Handler) getAllItems(c *gin.Context) {
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
	items, err := h.services.TodoItem.GetAll(userId, goalId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}
func (h *Handler) getItemById(c *gin.Context) {
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
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Item Id")
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId, goalId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, item)
	}
}
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Item id")
		return
	}
	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Goal id")
		return
	}

	var input todo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoItem.Update(userId, itemId, goalId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid User id")
		return
	}
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Item id")
		return
	}
	goalId, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid Goal id")
		return
	}
	err = h.services.TodoItem.Delete(userId, itemId, goalId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}
