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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email"})
		return
	}
	exists, err := h.services.UserExists(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User verification failed"})
		return
	}

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not exist"})
		return
	}
	code := h.services.Authorization.GenerateCode()

	if err := h.services.SendCodeToEmail(input.Email, code); err != nil {
		log.Println("error while sending code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send code"})
		return
	}

	mu.Lock()
	verificationCodes[input.Email] = VerificationCode{
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Email:     input.Email,
	}
	mu.Unlock()

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) verifyResetCode(c *gin.Context) {

	storedCode, err := checkCode(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if storedCode.IsVerified {
		c.JSON(http.StatusOK, statusResponse{"ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mail not confirmed"})
		return
	}
}

func (h *Handler) resetPassword(c *gin.Context) {
	var input struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect input"})
		return
	}
	err := validatePassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	storedCode := verificationCodes[input.Email]

	if !storedCode.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mail not confirmed"})
		return
	}

	err = h.services.Authorization.ResetPassword(input.Email, input.NewPassword)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	mu.Lock()
	delete(verificationCodes, input.Email)
	mu.Unlock()

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
