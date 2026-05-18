package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			var validationErr validator.ValidationErrors
			if errors.As(err.Err, &validationErr) {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "invalid input",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})

		}

	}
}
