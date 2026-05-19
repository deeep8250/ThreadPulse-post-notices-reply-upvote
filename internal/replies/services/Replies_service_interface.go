package services

import (
	"github.com/threadpulse/models"
)

type RepliesServiceInterface interface {
	CreateRepliesService(reply models.Replies) error
	GetAllRepliessService(postID int, limit, page int) ([]models.Replies, int, error)
	UpdateReplyService(UpdatedReply models.Replies) error
	DeleteReplyService(replyId, userID int) error
}
