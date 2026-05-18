package mock

import (
	"context"

	"github.com/threadpulse/models"
)

type ServiceMockThreads struct {
	CreateThreadFunc         func(userID int, input models.CreateThread) error
	GetAllThreadsFunc        func(page, limit int) ([]models.CreateThread, int, error)
	GetThreadByIdFunc        func(ThreadID int) (*models.CreateThread, error)
	UpdateThreadFunc         func(threadID, userID int, input models.UpdateThread) error
	DeleteThreadFunc         func(ThreadID, UserID int) error
	GetHotThreadsServiceFunc func(c context.Context, limit int) ([]models.HotThread, error)
}

func (s *ServiceMockThreads) CreateThread(userID int, input models.CreateThread) error {
	return s.CreateThreadFunc(userID, input)
}

func (s *ServiceMockThreads) GetAllThreads(page, limit int) ([]models.CreateThread, int, error) {
	return s.GetAllThreadsFunc(page, limit)
}
func (s *ServiceMockThreads) GetThreadById(ThreadID int) (*models.CreateThread, error) {
	return s.GetThreadByIdFunc(ThreadID)
}
func (s *ServiceMockThreads) UpdateThread(threadID, userID int, input models.UpdateThread) error {
	return s.UpdateThreadFunc(threadID, userID, input)
}
func (s *ServiceMockThreads) DeleteThread(ThreadID, UserID int) error {
	return s.DeleteThreadFunc(ThreadID, UserID)
}
func (s *ServiceMockThreads) GetHotThreadsService(c context.Context, limit int) ([]models.HotThread, error) {
	return s.GetHotThreadsServiceFunc(c, limit)
}
