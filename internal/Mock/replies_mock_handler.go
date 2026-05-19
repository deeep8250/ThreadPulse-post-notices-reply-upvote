package mock

import "github.com/threadpulse/models"

type MockingReplies struct {
	CreateRepliesServiceFunc  func(reply models.Replies) error
	GetAllRepliessServiceFunc func(postID int, limit, page int) ([]models.Replies, int, error)
	UpdateReplyServiceFunc    func(UpdatedReply models.Replies) error
	DeleteReplyServiceFunc    func(replyId, userID int) error
}

func (r *MockingReplies) CreateRepliesService(reply models.Replies) error {
	return r.CreateRepliesServiceFunc(reply)
}
func (r *MockingReplies) GetAllRepliessService(postID int, limit, page int) ([]models.Replies, int, error) {
	return r.GetAllRepliessServiceFunc(postID, limit, page)
}
func (r *MockingReplies) UpdateReplyService(UpdatedReply models.Replies) error {
	return r.UpdateReplyServiceFunc(UpdatedReply)
}
func (r *MockingReplies) DeleteReplyService(replyId, userID int) error {
	return r.DeleteReplyServiceFunc(replyId, userID)
}
