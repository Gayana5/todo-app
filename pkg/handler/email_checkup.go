package handler

import (
	"fmt"
	"net/smtp"
	"time"
)

const (
	SMTP_SERVER   = "smtp.mail.ru"
	SMTP_PORT     = "587"
	SMTP_USERNAME = "whattodo.confirm@mail.ru"
	SMTP_PASSWORD = "gPHcX4wNZbHSYdZsi3WX"
)

type VerificationCode struct {
	Code      string
	ExpiresAt time.Time
}

var verificationCodes = make(map[string]VerificationCode)

func sendCodeToEmail(email, code string) error {
	auth := smtp.PlainAuth("", SMTP_USERNAME, SMTP_PASSWORD, SMTP_SERVER)
	to := []string{email}
	subject := "Код подтверждения"
	body := fmt.Sprintf("Ваш код подтверждения: %s", code)

	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	return smtp.SendMail(SMTP_SERVER+":"+SMTP_PORT, auth, SMTP_USERNAME, to, msg)
}
