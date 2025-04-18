package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/Gayana5/todo-app"
	"github.com/Gayana5/todo-app/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

const (
	SALT        =
	SIGNING_KEY =
	tokenTTL    = 720 * time.Hour
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
