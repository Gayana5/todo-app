package handler

import (
	"github.com/Gayana5/todo-app"
	_ "github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var mu sync.Mutex

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	/*
			exists, err := h.repo.UserExists(input.Email)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки пользователя"})
				return
			}
			if exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с такой почтой уже существует"})
				return
			}

			code := h.services.Authorization.GenerateCode()

			if err := sendCodeToEmail(input.Email, code); err != nil {
				log.Println("Ошибка отправки кода:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить код"})
				return
			}

			mu.Lock()
			verificationCodes[input.Email] = VerificationCode{
				Code:      code,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"message": "Код подтверждения отправлен на вашу почту."})
	*/
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) verifyCode(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	mu.Lock()
	storedCode, exists := verificationCodes[input.Email]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Код не существует"})
		return
	}
	if storedCode.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Срок действия кода истек"})
		return
	}

	if storedCode.Code != input.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный код"})
		return
	}

	mu.Lock()
	delete(verificationCodes, input.Email)
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Регистрация успешна"})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()) // Error 400. Пользователь предоставил некорректные данные.
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
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

func (h *Handler) getInfo(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь не найден."})
		return
	}
	user, err := h.services.Authorization.GetInfo(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"first_name":  user.FirstName,
		"second_name": user.SecondName,
		"email":       user.Email,
	})
}
