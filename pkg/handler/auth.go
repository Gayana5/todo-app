package handler

import (
	"github.com/Gayana5/todo-app"
	_ "github.com/Gayana5/todo-app"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

type VerificationCode struct {
	Code       string
	ExpiresAt  time.Time
	UserData   todo.User
	Email      string
	IsVerified bool
}

var (
	verificationCodes = make(map[string]VerificationCode)
	mu                sync.Mutex
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	err := validatePassword(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	exists, err := h.services.UserExists(input.Email)
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
		UserData:  input,
		Email:     input.Email,
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"code": code})
}

func (h *Handler) verifyRegistrationCode(c *gin.Context) {
	storedCode, err := checkCode(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	mu.Lock()
	userData := storedCode.UserData
	delete(verificationCodes, storedCode.Email)
	mu.Unlock()

	_, err = h.services.Authorization.CreateUser(userData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(userData.Email, userData.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // Error 500. Внутренняя ошибка на сервере.
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
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
