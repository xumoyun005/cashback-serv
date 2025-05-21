package service

import (
	constants "cashback-serv/const"
	"cashback-serv/internal/interfaces"
	"cashback-serv/internal/queue"
	"cashback-serv/models"
	"errors"
	"fmt"
	"time"
)

type CashbackRepository interface {
	GetCashbackByUserID(userID int64, fromDate, toDate string) (*models.Cashback, error)
	CreateCashback(cashback *models.Cashback) error
	UpdateCashbackAmount(id int64, amount float64) error
	CreateCashbackHistory(history *models.CashbackHistory) error
	GetCashbackHistoryByUserID(turonUserID int64, fromDate, toDate string, pagination *models.Pagination) ([]models.CashbackHistory, error)
}

type CashbackService struct {
	repo          CashbackRepository
	queue         *queue.CashbackQueue
	sourceService core.SourceFinderCreator
}

func NewCashbackService(repo CashbackRepository, sourceService core.SourceFinderCreator) *CashbackService {
	return &CashbackService{
		repo:          repo,
		queue:         queue.NewCashbackQueue(repo, sourceService),
		sourceService: sourceService,
	}
}

func (s *CashbackService) IncreaseCashback(req *models.CashbackRequest) error {
	if req.TuronUserID != 0 && req.CineramaUserID != 0 {
		return errors.New("only one of turon_user_id or cinerama_user_id should be provided")
	}

	source, err := s.sourceService.FindSourceOrCreate(req.TuronUserID, req.CineramaUserID, req.HostIP)
	if err != nil {
		return fmt.Errorf("failed to determine source: %w", err)
	}

	return s.queue.Enqueue(constants.Increase, req, source.ID)
}

func (s *CashbackService) DecreaseCashback(req *models.CashbackRequest) error {
	if req.TuronUserID != 0 && req.CineramaUserID != 0 {
		return errors.New("only one of turon_user_id or cinerama_user_id should be provided")
	}

	source, err := s.sourceService.FindSourceOrCreate(req.TuronUserID, req.CineramaUserID, req.HostIP)
	if err != nil {
		return fmt.Errorf("failed to determine source: %w", err)
	}

	return s.queue.Enqueue(constants.Decrease, req, source.ID)
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
