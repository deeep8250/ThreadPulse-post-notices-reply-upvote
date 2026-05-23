package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/threadpulse/internal/replies/services"
	"github.com/threadpulse/models"
)

type RepliesHandler struct {
	service service.RepliesServiceInterface
}

func NewRepliesHandler(serv service.RepliesServiceInterface) *RepliesHandler {
	return &RepliesHandler{
		service: serv,
	}
}

func (h *RepliesHandler) CreateRepliesHandler(c *gin.Context) {

	postID := c.Param("id")
	postIDint, err := strconv.Atoi(postID)
	if err != nil || postIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	var reply models.Replies
	err = c.ShouldBindJSON(&reply)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
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

	reply.PostID = postIDint
	reply.UserID = userID.(int)

	err = h.service.CreateRepliesService(reply)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "reply added to post's comment section",
	})

}

func (h *RepliesHandler) GetAllRepliesHandler(c *gin.Context) {

	PostId := c.Param("id")
	PostIdInt, err := strconv.Atoi(PostId)
	if err != nil || PostIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	limit := c.DefaultQuery("limit", "5")

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}
	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	replies, count, err := h.service.GetAllRepliessService(PostIdInt, limitInt, pageInt)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"postID":  PostId,
		"replies": replies,
		"count":   count,
	})

}

func (h *RepliesHandler) UpdateRepliesHandler(c *gin.Context) {
	ReplyID := c.Param("id")
	ReplyIdInt, err := strconv.Atoi(ReplyID)
	if err != nil || ReplyIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid url parameter",
		})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	var replyUpdated models.Replies
	err = c.ShouldBindJSON(&replyUpdated)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	replyUpdated.Id = ReplyIdInt
	replyUpdated.UserID = userID.(int)

	err = h.service.UpdateReplyService(replyUpdated)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "reply updated",
	})

}

func (h *RepliesHandler) DeleteReplyHandler(c *gin.Context) {
	replyID := c.Param("id")
	replyIdInt, err := strconv.Atoi(replyID)
	if err != nil || replyIdInt <= 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reply id",
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

	err = h.service.DeleteReplyService(replyIdInt, userID.(int))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "deleted the reply",
	})

}
