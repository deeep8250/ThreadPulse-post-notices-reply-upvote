package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	mock "github.com/threadpulse/internal/Mock"
	"github.com/threadpulse/internal/middleware"

	"github.com/threadpulse/models"
)

func TestCreateThreadHandler(t *testing.T) {
	os.Setenv("JWT_SECRET", "your_test_secret")
	tests := []struct {
		TestName           string
		UserInput          models.CreateThread
		Token              string
		ServiceError       error
		ExpectedStatusCode int
	}{
		{
			TestName: "success",
			UserInput: models.CreateThread{
				UserID:  1,
				Title:   "hw",
				Content: "hello world i got a placement yoho !",
			},
			Token:              "Bearer " + GenerateTestToken(),
			ServiceError:       nil,
			ExpectedStatusCode: 201,
		},

		{
			TestName: "service error",
			UserInput: models.CreateThread{
				UserID:  1,
				Title:   "hw",
				Content: "hello world i got a placement yoho !",
			},
			Token:        "Bearer " + GenerateTestToken(),
			ServiceError: errors.New("something failed"), ExpectedStatusCode: 500},

		{
			TestName:           "invalid input",
			UserInput:          models.CreateThread{},
			Token:              "Bearer " + GenerateTestToken(),
			ServiceError:       nil,
			ExpectedStatusCode: 400,
		},
		{
			TestName:           "no token",
			UserInput:          models.CreateThread{},
			Token:              "",
			ServiceError:       nil,
			ExpectedStatusCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {

			mockingService := mock.ServiceMockThreads{
				CreateThreadFunc: func(userID int, input models.CreateThread) error {
					return tt.ServiceError
				},
			}

			handler := NewThreadHandler(&mockingService)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.Use(middleware.Miiddleware())
			r.POST("/thread", handler.CreateThreadHandler)
			bodyBytes, _ := json.Marshal(tt.UserInput)

			req, _ := http.NewRequest(http.MethodPost, "/thread", bytes.NewBufferString(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")
			if tt.Token != "" {
				req.Header.Set("Authorization", tt.Token) // ← attach it
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.ExpectedStatusCode {
				t.Errorf("Expected this %d got this %d", tt.ExpectedStatusCode, w.Code)
			}
		})
	}

}

func GenerateTestToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": 1,
	})
	tokenString, _ := token.SignedString([]byte("your_test_secret"))
	return tokenString
}
