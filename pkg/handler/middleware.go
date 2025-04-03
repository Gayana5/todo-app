package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user_id"
)

func (h *Handler) userIdentity(c *gin.Context) {

	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "No authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}
	userId, err := h.services.Authorization.ParseToken(headerParts[1])

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("User Id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("User Id is of invalid type")
	}
	return idInt, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password is too short")
	}

	hasLower := false
	hasUpper := false
	hasDigit := false

	for _, ch := range password {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		}
	}

	if !(hasLower && hasUpper && hasDigit) {
		return errors.New("invalid password")
	}

	return nil
}
func checkCode(c *gin.Context) (VerificationCode, error) {
	var input struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BindJSON(&input); err != nil {

		return VerificationCode{}, errors.New("incorrect input")
	}

	mu.Lock()
	storedCode, exists := verificationCodes[input.Email]
	mu.Unlock()

	if !exists {
		return VerificationCode{}, errors.New("code not found")
	}
	if storedCode.ExpiresAt.Before(time.Now()) {
		return VerificationCode{}, errors.New("code expired")
	}

	if storedCode.Code != input.Code {
		return VerificationCode{}, errors.New("incorrect code")
	}

	storedCode.IsVerified = true

	mu.Lock()
	verificationCodes[input.Email] = storedCode
	mu.Unlock()

	return storedCode, nil
}
