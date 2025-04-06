package handler

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
)

const (
	SMPT_HOST     = "smtp.mail.ru"
	SMTP_PORT     = 587
	SMTP_USERNAME =
	SMTP_PASSWORD =
)

func sendCodeToEmail(to string, code string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", SMTP_USERNAME)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "WhatToDo - Confirmation Code")
	m.SetBody("text/plain", fmt.Sprintf("Your onetime verification code: %s", code))

	d := gomail.NewDialer(SMPT_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Игнорируем проверку сертификата

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}
