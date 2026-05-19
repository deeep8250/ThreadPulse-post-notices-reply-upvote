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
			r.Use(middleware.ErrorHandler(), middleware.Miiddleware())

			r.POST("/thread/:id", handler.CreateThreadHandler)
			bodyBytes, _ := json.Marshal(tt.UserInput)

			req, _ := http.NewRequest(http.MethodPost, "/thread/2", bytes.NewBufferString(string(bodyBytes)))
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

func TestSUpdateThreadHandler(t *testing.T) {
	os.Setenv("JWT_SECRET", "your_test_secret")
	tests := []struct {
		TestName     string
		Token        string
		UserInput    models.UpdateThread
		ServiceErr   error
		ExepctedCode int
	}{
		{
			TestName: "Success",
			Token:    "Bearer " + GenerateTestToken(),
			UserInput: models.UpdateThread{
				Title:   "test",
				Content: "testing ",
			},
			ServiceErr:   nil,
			ExepctedCode: 200,
		},
		{
			TestName: "invalid input",
			Token:    "Bearer " + GenerateTestToken(),
			UserInput: models.UpdateThread{
				Title:   "",
				Content: "",
			},
			ServiceErr:   nil,
			ExepctedCode: 500,
		},
		{
			TestName: "unauthorized user",
			Token:    "",
			UserInput: models.UpdateThread{
				Title:   "test",
				Content: "testing ",
			},
			ServiceErr:   nil,
			ExepctedCode: 401,
		},
		{
			TestName:     "service error",
			Token:        "Bearer " + GenerateTestToken(),
			UserInput:    models.UpdateThread{Title: "test", Content: "testing"},
			ServiceErr:   errors.New("something failed"),
			ExepctedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			mocking := mock.ServiceMockThreads{
				UpdateThreadFunc: func(threadID, userID int, input models.UpdateThread) error {
					return tt.ServiceErr
				},
			}

			handler := NewThreadHandler(&mocking)
			r := gin.Default()
			r.Use(middleware.ErrorHandler(), middleware.Miiddleware())
			r.POST("/thread/:id", handler.UpdateThreadHandler)

			bodyBytes, _ := json.Marshal(tt.UserInput)
			req, _ := http.NewRequest(http.MethodPost, "/thread/1", bytes.NewBufferString(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")
			if tt.Token != "" {
				req.Header.Set("Authorization", tt.Token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.ExepctedCode {
				t.Errorf("Expected this %d got this %d", tt.ExepctedCode, w.Code)
			}

		})
	}

}

func TestDeleteThreadHandler(t *testing.T) {
	os.Setenv("JWT_SECRET", "your_test_secret")
	tests := []struct {
		TestName     string
		Token        string
		ThreadID     string
		ServiceErr   error
		ExpectedCode int
	}{
		{
			TestName:     "success",
			Token:        "Bearer " + GenerateTestToken(),
			ThreadID:     "1",
			ServiceErr:   nil,
			ExpectedCode: 200,
		},
		{
			TestName:     "invalid thread id",
			Token:        "Bearer " + GenerateTestToken(),
			ThreadID:     "abc",
			ServiceErr:   nil,
			ExpectedCode: 500, // or 400 if your ErrorHandler returns BadRequest
		},
		{
			TestName:     "unauthorized user",
			Token:        "",
			ThreadID:     "1",
			ServiceErr:   nil,
			ExpectedCode: 401,
		},
		{
			TestName:     "service error",
			Token:        "Bearer " + GenerateTestToken(),
			ThreadID:     "1",
			ServiceErr:   errors.New("something failed"),
			ExpectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {

			mocking := mock.ServiceMockThreads{
				DeleteThreadFunc: func(ThreadID, UserID int) error {
					return tt.ServiceErr
				},
			}

			handler := NewThreadHandler(&mocking)
			r := gin.Default()
			r.Use(middleware.ErrorHandler(), middleware.Miiddleware())
			r.DELETE("/thread/:id", handler.DeleteThreadHandler)

			req, _ := http.NewRequest(http.MethodDelete, "/thread/"+tt.ThreadID, nil)
			req.Header.Set("Content-Type", "appplication/json")
			if tt.Token != "" {
				req.Header.Set("Authorization", tt.Token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.ExpectedCode {
				t.Errorf("expected %d got %d", tt.ExpectedCode, w.Code)
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
