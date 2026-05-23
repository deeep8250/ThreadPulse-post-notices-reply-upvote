package services

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/threadpulse/internal/config"
	"github.com/threadpulse/internal/upvotes/repositories"
)

type UpvoteService struct {
	repo        *repositories.UpvotesRepository
	worker      *repositories.UpvoteWorker
	redisClient *redis.Client
}

func NewUpvoteService(Repo *repositories.UpvotesRepository, Worker *repositories.UpvoteWorker) *UpvoteService {
	return &UpvoteService{
		repo:        Repo,
		worker:      Worker,
		redisClient: config.RedisClient,
	}
}

func (s *UpvoteService) SubmitUpvote(postID, userID int, c context.Context) error {
	exists, err := s.repo.CheckUpvote(postID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already upvoted")
	}
	s.worker.Submit(postID, userID)
	r := s.redisClient.Del(c, "hot_threads")
	log.Println("cache deleted keys:", r.Val())
	return nil

}

func (s *UpvoteService) GetUpvotes(postID int) (int, error) {

	// check thread id is available or not

	upvotes, err := s.repo.GetUpvotes(postID)
	if err != nil {
		return 0, err
	}
	if upvotes == 0 {
		return 0, errors.New("invalid post id")
	}
	return upvotes, nil

}
