package service

import (
	constants "cashback-serv/const"
	"cashback-serv/models"
	"cashback-serv/internal/queue"
	"cashback-serv/internal/repository"
)

type CashbackService struct {
	repo  *repository.CashbackRepository
	queue *queue.CashbackQueue
}

func NewCashbackService(repo *repository.CashbackRepository) *CashbackService {
	return &CashbackService{
		repo:  repo,
		queue: queue.NewCashbackQueue(repo),
	}
}

func (s *CashbackService) IncreaseCashback(req *models.CashbackRequest) error {
	return s.queue.Enqueue(constants.Increase, req)
}

func (s *CashbackService) DecreaseCashback(req *models.CashbackRequest) error {
	return s.queue.Enqueue(constants.Decrease, req)
}

func (s *CashbackService) GetCashbackByUserID(turonUserID int64) (*models.Cashback, error) {
	return s.repo.GetCashbackByUserID(turonUserID)
}

func (s *CashbackService) GetCashbackHistoryByUserID(turonUserID int64) ([]models.CashbackHistory, error) {
	return s.repo.GetCashbackHistoryByUserID(turonUserID)
}
