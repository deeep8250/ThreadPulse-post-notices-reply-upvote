package services

import (
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/threadpulse/internal/config"
	"github.com/threadpulse/internal/threads/repository"
	"github.com/threadpulse/models"
	"golang.org/x/net/context"
)

type ThreadsService struct {
	repo        *repository.ThreadsRepo
	redisClient *redis.Client
}

func NewThreadsService(repo *repository.ThreadsRepo) *ThreadsService {
	return &ThreadsService{
		repo:        repo,
		redisClient: config.RedisClient,
	}

}

func (s *ThreadsService) CreateThread(userID int, input models.CreateThread) error {

	input.UserID = userID

	err := s.repo.CreateThreads(&input)
	if err != nil {
		return err
	}
	return nil
}

func (s *ThreadsService) GetAllThreads(page, limit int) ([]models.CreateThread, int, error) {
	offset := (page - 1) * limit

	threads, count, err := s.repo.GetAllThreads(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return *threads, count, nil

}

func (s *ThreadsService) GetThreadById(ThreadID int) (*models.CreateThread, error) {
	thread, err := s.repo.GetThreadByID(ThreadID)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (s *ThreadsService) UpdateThread(threadID, userID int, input models.UpdateThread) error {
	err := s.repo.UpdateThread(threadID, userID, input)
	if err != nil {
		return err
	}
	return nil
}

func (s *ThreadsService) DeleteThread(ThreadID, UserID int) error {
	err := s.repo.DeleteThread(ThreadID, UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ThreadsService) GetHotThreadsService(c context.Context, limit int) ([]models.HotThread, error) {

	val, err := s.redisClient.Get(c, "hot_threads").Result()
	if err == redis.Nil {
		hotThreads, err := s.repo.GetHotThread(limit)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(hotThreads)
		if err != nil {
			return nil, err
		}

		s.redisClient.Set(c, "hot_threads", data, 5*time.Minute)
		return hotThreads, nil

	} else if err != nil {
		return nil, err
	} else {
		var hotThreads []models.HotThread
		json.Unmarshal([]byte(val), &hotThreads)
		return hotThreads, nil
	}

}
