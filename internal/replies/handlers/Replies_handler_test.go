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
func TestUpdateReplyService(t *testing.T) {
	os.Setenv("JWT_SECRET", "your_test_secret")
	tests := []struct {
		TestName     string
		Token        string
		ThreadID     string
		UserInput    models.Replies
		ServiceErr   error
		ExpectedCode int
	}{
		{
			TestName: "success",
			Token:    "Bearer " + GenerateTestToken(),
			ThreadID: "1",
			UserInput: models.Replies{

				Reply: "hello world",
			},
			ServiceErr:   nil,
			ExpectedCode: 200,
		},
		{
			TestName: "invalid input",
			ThreadID: "1",
			Token:    "Bearer " + GenerateTestToken(),
			UserInput: models.Replies{

				PostID: 1,
				Reply:  "",
			},
			ServiceErr:   nil,
			ExpectedCode: 400,
		},
		{
			TestName: "unauthorized user",
			ThreadID: "1",
			Token:    "",
			UserInput: models.Replies{
				PostID: 1,
				Reply:  "hello world",
			},
			ServiceErr:   nil,
			ExpectedCode: 401,
		},
		{
			TestName: "service error",
			ThreadID: "1",
			Token:    "Bearer " + GenerateTestToken(),
			UserInput: models.Replies{
				PostID: 1,
				Reply:  "hello world",
			},

			ServiceErr:   errors.New("something wents wrong"),
			ExpectedCode: 500,
		},
		{
			TestName: "invalid post id",
			ThreadID: "0",
			Token:    "Bearer " + GenerateTestToken(),
			UserInput: models.Replies{
				PostID: 0,
				Reply:  "hello world",
			},

			ServiceErr:   nil,
			ExpectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			mocking := mock.MockingReplies{
				UpdateReplyServiceFunc: func(UpdatedReply models.Replies) error {
					return tt.ServiceErr
				},
			}
			handler := NewRepliesHandler(&mocking)
			r := gin.Default()
			r.Use(middleware.ErrorHandler(), middleware.Miiddleware())
			r.PATCH("/replies/:id", handler.UpdateRepliesHandler)
			bodyBytes, _ := json.Marshal(tt.UserInput)
			req, _ := http.NewRequest(http.MethodPatch, "/replies/"+tt.ThreadID, bytes.NewBufferString(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")
			if tt.Token != "" {
				req.Header.Set("Authorization", tt.Token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.ExpectedCode {
				t.Errorf("Expected %d got %d ", tt.ExpectedCode, w.Code)
			}

		})
	}

}

func TestDeleteReplyService(t *testing.T) {
	os.Setenv("JWT_SECRET", "your_test_secret")
	tests := []struct {
		TestName     string
		Token        string
		ReplyID      string
		ServiceErr   error
		ExpectedCode int
	}{
		{
			TestName:     "success",
			Token:        "Bearer " + GenerateTestToken(),
			ReplyID:      "1",
			ServiceErr:   nil,
			ExpectedCode: 200,
		},
		{
			TestName: "invalid input",
			ReplyID:  "0",
			Token:    "Bearer " + GenerateTestToken(),

			ServiceErr:   nil,
			ExpectedCode: 400,
		},
		{
			TestName: "unauthorized user",
			ReplyID:  "1",
			Token:    "",

			ServiceErr:   nil,
			ExpectedCode: 401,
		},
		{
			TestName: "service error",
			ReplyID:  "1",
			Token:    "Bearer " + GenerateTestToken(),

			ServiceErr:   errors.New("something wents wrong"),
			ExpectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			mocking := mock.MockingReplies{
				DeleteReplyServiceFunc: func(replyId, userID int) error {
					return tt.ServiceErr
				},
			}
			handler := NewRepliesHandler(&mocking)
			r := gin.Default()
			r.Use(middleware.ErrorHandler(), middleware.Miiddleware())
			r.DELETE("/replies/:id", handler.DeleteReplyHandler)

			req, _ := http.NewRequest(http.MethodDelete, "/replies/"+tt.ReplyID, nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.Token != "" {
				req.Header.Set("Authorization", tt.Token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.ExpectedCode {
				t.Errorf("Expected %d got %d ", tt.ExpectedCode, w.Code)
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
