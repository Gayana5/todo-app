package service

import (
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

const (
	SALT        = "ncuewfr53567njwejk95"
	SIGNING_KEY = "hfwoiujr8420#fiopsrUHfewijfHe"
	tokenTTL    = 720 * time.Hour

	SMPT_HOST     = "smtp.mail.ru"
	SMTP_PORT     = 587
	SMTP_USERNAME = "whattodo.confirm@mail.ru"
	SMTP_PASSWORD = "gPHcX4wNZbHSYdZsi3WX"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}
func (s *AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repo.GetUser(email, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(SIGNING_KEY))
}
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(SIGNING_KEY), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}
	return claims.UserId, nil
}
func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(SALT)))
}

func (s *AuthService) GenerateCode() string {

	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func (s *AuthService) GetInfo(id int) (todo.User, error) {
	return s.repo.GetInfo(id)
}
func (s *AuthService) UpdateInfo(userId int, input todo.UpdateUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateInfo(userId, input)
}
func (s *AuthService) UserExists(email string) (bool, error) {
	return s.repo.UserExists(email)
}
func (s *AuthService) ResetPassword(email, password string) error {
	password = s.generatePasswordHash(password)
	return s.repo.ResetPassword(email, password)
}

func (s *AuthService) SendCodeToEmail(to string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", SMTP_USERNAME)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "WhatToDo - Confirmation Code")
	m.SetBody("text/plain", fmt.Sprintf("Your onetime verification code: %s", code))

	d := gomail.NewDialer(SMPT_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}
