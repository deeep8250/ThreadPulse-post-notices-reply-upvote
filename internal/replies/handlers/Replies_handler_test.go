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

func TestCreateReplies(t *testing.T) {
	os.Setenv("JWT_SECRET", "your_test_secret")
	tests := []struct {
		TestName     string
		UserInput    models.Replies
		Token        string
		ExpectedCode int
		ServiceErr   error
	}{
		{
			TestName: "success",
			UserInput: models.Replies{

				Reply: "helllo world",
			},
			Token:        "Bearer " + GenerateTestToken(),
			ExpectedCode: 201,
			ServiceErr:   nil,
		},
		{
			TestName: "post not found",
			UserInput: models.Replies{

				Reply: "helllo world",
			},
			Token:        "Bearer " + GenerateTestToken(),
			ExpectedCode: 500,
			ServiceErr:   errors.New("something failed"),
		},

		{

			TestName: "Invalid input",
			UserInput: models.Replies{

				Reply: "",
			},
			Token:        "Bearer " + GenerateTestToken(),
			ExpectedCode: 400,
			ServiceErr:   nil,
		},

		{

			TestName: "no token",
			UserInput: models.Replies{

				Reply: "",
			},
			Token:        "",
			ExpectedCode: 401,
			ServiceErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			mockingService := &mock.MockingReplies{
				CreateRepliesServiceFunc: func(reply models.Replies) error {
					return tt.ServiceErr
				},
			}
			handler := NewRepliesHandler(mockingService)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.Use(middleware.Miiddleware())
			r.POST("/thread/:id/reply", handler.CreateRepliesHandler)
			bodyBytes, _ := json.Marshal(tt.UserInput)
			req, _ := http.NewRequest(http.MethodPost, "/thread/1/reply", bytes.NewBufferString(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")
			if tt.Token != "" {
				req.Header.Set("Authorization", tt.Token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.ExpectedCode {
				t.Errorf("Expected this %d got this %d", tt.ExpectedCode, w.Code)
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
