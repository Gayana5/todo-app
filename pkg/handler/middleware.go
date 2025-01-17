package handler

import (
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
