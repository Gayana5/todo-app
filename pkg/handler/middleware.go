package handler

import (
	"errors"
	"github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	// Получаем header авторизации.
	header := c.GetHeader(authorizationHeader)
	// Валидируем, что header не пустой.
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "No authorization header") // Error 401.
		return
	}
	// Валидируем, что header корректный.
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header") // Error 401. Пользователь не авторизирован
		return
	}
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	// Валидируем на наличие других возможных ошибок.
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "User id not found")
		return 0, errors.New("User id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "User id is of incorrect type")
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}
func validateUser(user todo.User) error {
	if !nameRegex.MatchString(user.FirstName) {
		return errors.New("first_name должно содержать 3-20 символов, включая буквы, цифры, _ и -")
	}
	if !nameRegex.MatchString(user.SecondName) {
		return errors.New("second_name должно содержать 3-20 символов, включая буквы, цифры, _ и -")
	}
	if !passwordRegex.MatchString(user.Password) {
		return errors.New("пароль должен содержать 5-30 символов, включая буквы, цифры, _ и -")
	}
	return nil
}
