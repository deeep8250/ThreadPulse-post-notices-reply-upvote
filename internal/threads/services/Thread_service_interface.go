package services

import (
	"github.com/threadpulse/models"
	"golang.org/x/net/context"
)

type ServiceInterface interface {
	CreateThread(userID int, input models.CreateThread) error
	GetAllThreads(page, limit int) ([]models.CreateThread, int, error)
	GetThreadById(ThreadID int) (*models.CreateThread, error)
	UpdateThread(threadID, userID int, input models.UpdateThread) error
	DeleteThread(ThreadID, UserID int) error
	GetHotThreadsService(c context.Context, limit int) ([]models.HotThread, error)
}
