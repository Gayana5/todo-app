package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Gayana5/todo-app/pkg/service"
	"github.com/Gayana5/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type signUpRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func TestHandler_signUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthorization(ctrl)

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	router := gin.Default()
	router.POST("/sign-up", handler.signUp)

	t.Run("Success", func(t *testing.T) {
		mockAuth.EXPECT().UserExists("test@example.com").Return(false, nil)
		mockAuth.EXPECT().GenerateCode().Return("1234")
		mockAuth.EXPECT().SendCodeToEmail("test@example.com", "1234").Return(nil)

		req := signUpRequest{
			FirstName:  "Test",
			SecondName: "Test",
			Email:      "test@example.com",
			Password:   "Qwerty123",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, reqHttp)

		require.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})

	t.Run("Incorrect_Password", func(t *testing.T) {
		mockAuth.EXPECT().UserExists("test@example.com").Return(false, nil)
		mockAuth.EXPECT().GenerateCode().Return("123456")

		req := signUpRequest{
			FirstName:  "Test",
			SecondName: "Test",
			Email:      "test@example.com",
			Password:   "123", // слишком короткий пароль
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message":"incorrect password"}`, w.Body.String())
	})

	t.Run("Failed_to_send_code", func(t *testing.T) {
		mockAuth.EXPECT().SendCodeToEmail("test@example.com", "123456").Return(errors.New("smtp error"))

		req := signUpRequest{
			FirstName:  "Test",
			SecondName: "Test",
			Email:      "test@example.com",
			Password:   "Qwerty123",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message":"failed to send code"}`, w.Body.String())
	})
}
