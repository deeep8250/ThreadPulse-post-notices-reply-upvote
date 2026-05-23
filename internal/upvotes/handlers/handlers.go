package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/threadpulse/internal/upvotes/services"
)

type UpvoteHandler struct {
	service *services.UpvoteService
}

func NewUpvoteHandler(Serv *services.UpvoteService) *UpvoteHandler {
	return &UpvoteHandler{
		service: Serv,
	}
}

func (h *UpvoteHandler) Upvote(c *gin.Context) {
	postIDstring := c.Param("id")
	postID, err := strconv.Atoi(postIDstring)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})
		return
	}

	err = h.service.SubmitUpvote(postID, userID.(int), c.Request.Context())
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "upvoted",
	})

}

func (h *UpvoteHandler) GetAllUpvotes(c *gin.Context) {
	postIDstring := c.Param("id")
	postID, err := strconv.Atoi(postIDstring)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	countVotes, err := h.service.GetUpvotes(postID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"upvotes": countVotes,
	})
}
