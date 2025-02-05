package handler

import (
	"github.com/Gayana5/todo-app"
	_ "github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if err := validateUser(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()) // Error 400. Пользователь предоставил некорректные данные.
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // Error 500. Внутренняя ошибка на сервере.
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

var (
	nameRegex     = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{5,30}$`)
)
