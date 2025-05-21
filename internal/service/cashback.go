package service

import (
	constants "cashback-serv/const"
	"cashback-serv/internal/queue"
	"cashback-serv/internal/repository"
	"cashback-serv/models"
	"errors"
	"time"
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

func (s *CashbackService) GetCashbackByUserID(turonUserID int64, fromDate, toDate string) (*models.Cashback, error) {
	if err := s.validateDates(fromDate, toDate); err != nil {
		return nil, err
	}
	return s.repo.GetCashbackByUserID(turonUserID, fromDate, toDate)
}

func (s *CashbackService) GetCashbackHistoryByUserID(turonUserID int64, fromDate, toDate string, pagination *models.Pagination) ([]models.CashbackHistory, error) {
	if err := s.validateDates(fromDate, toDate); err != nil {
		return nil, err
	}

	if err := s.validatePagination(pagination); err != nil {
		return nil, err
	}

	pagination.Calculate()
	return s.repo.GetCashbackHistoryByUserID(turonUserID, fromDate, toDate, pagination)
}

func (s *CashbackService) validateDates(fromDate, toDate string) error {
	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			return errors.New("invalid from_date format. Use YYYY-MM-DD")
		}
	}
	if toDate != "" {
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			return errors.New("invalid to_date format. Use YYYY-MM-DD")
		}
	}
	return nil
}

func (s *CashbackService) validatePagination(pagination *models.Pagination) error {
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 {
		pagination.PageSize = 10
	}
	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}
	return nil
}
