package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (h *Handler) forgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	exists, err := h.services.UserExists(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки пользователя"})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с такой почтой не существует"})
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
		Email:     input.Email,
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Код подтверждения отправлен на вашу почту."})
}

func (h *Handler) verifyResetCode(c *gin.Context) {

	storedCode, err := checkCode(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if storedCode.IsVerified {
		c.JSON(http.StatusOK, statusResponse{"ok"})
	} else {
		newErrorResponse(c, http.StatusBadRequest, "почта не подтверждена")
		return
	}
}

func (h *Handler) resetPassword(c *gin.Context) {
	var input struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	err := validatePassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	storedCode := verificationCodes[input.Email]

	if !storedCode.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Почта не подтверждена"})
		return
	}

	err = h.services.Authorization.ResetPassword(input.Email, input.NewPassword)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	mu.Lock()
	delete(verificationCodes, input.Email)
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "пароль успешно изменен"})
}
