package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mock "github.com/threadpulse/internal/Mock"
	"github.com/threadpulse/internal/middleware"
	"github.com/threadpulse/models"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		TestName           string
		UserInput          models.Register
		ServiceError       error
		ExpectedStatusCode int
	}{
		{TestName: "success", UserInput: models.Register{Name: "Deep mondal", Email: "Deep@gmail.com", Password: "Deep@123"}, ServiceError: nil, ExpectedStatusCode: 201},
		{TestName: "invalid input", UserInput: models.Register{Name: "Deep mondal", Email: "Deepgmail.com", Password: "Deep@123"}, ServiceError: errors.New("invalid request"), ExpectedStatusCode: 400},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {

			mockingService := mock.ServiceMock{
				RegisterFunc: func(registerInput models.Register) error {
					return tt.ServiceError
				},
			}

			handler := NewAuthHandler(&mockingService)

			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.POST("/register", handler.RegisterHandler)
			bodyBytes, _ := json.Marshal(tt.UserInput)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.ExpectedStatusCode {
				t.Errorf("Expected this %d got this %d", tt.ExpectedStatusCode, w.Code)
			}

		})
	}

}

func TestLogin(t *testing.T) {
	tests := []struct {
		TestName           string
		UserInput          models.Login
		ServiceError       error
		Token              string
		ExpectedStatusCode int
	}{
		{TestName: "success", UserInput: models.Login{Email: "Deep@gmail.com", Password: "Deep@123"}, ServiceError: nil, Token: "abcd", ExpectedStatusCode: 200},
		{TestName: "invalid input", UserInput: models.Login{Email: "Deepgmail.com", Password: "Deep@123"}, ServiceError: errors.New("invalid request"), Token: "", ExpectedStatusCode: 400},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			mockingService := mock.ServiceMock{
				LoginFunc: func(user models.Login) (string, error) {
					return tt.Token, tt.ServiceError
				},
			}

			handler := NewAuthHandler(&mockingService)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.POST("/login", handler.Login)
			bodyBytes, _ := json.Marshal(tt.UserInput)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.ExpectedStatusCode {
				t.Errorf("Expected this %d got this %d", tt.ExpectedStatusCode, w.Code)
			}
		})
	}
}
