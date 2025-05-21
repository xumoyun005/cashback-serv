package service

import (
	constants "cashback-serv/const"
	core "cashback-serv/internal/interfaces"
	"cashback-serv/internal/queue"
	"cashback-serv/models"
	"errors"
	"fmt"
	"time"
)

type CashbackRepository interface {
	GetCashbackByUserID(turonUserID int64) (*models.Cashback, error)
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

func (s *CashbackService) validateTuronUserID(turonUserID int64) error {
	if turonUserID == 0 {
		return errors.New("turon_user_id must be provided")
	}
	return nil
}

func (s *CashbackService) IncreaseCashback(req *models.CashbackRequest) error {
	if err := s.validateTuronUserID(req.TuronUserID); err != nil {
		return err
	}

	source, err := s.sourceService.FindSourceOrCreate(req.TuronUserID, req.HostIP)
	if err != nil {
		return fmt.Errorf("failed to determine source: %w", err)
	}

	return s.queue.Enqueue(constants.Increase, req, source.ID)
}

func (s *CashbackService) DecreaseCashback(req *models.CashbackRequest) error {
	if err := s.validateTuronUserID(req.TuronUserID); err != nil {
		return err
	}

	source, err := s.sourceService.FindSourceOrCreate(req.TuronUserID, req.HostIP)
	if err != nil {
		return fmt.Errorf("failed to determine source: %w", err)
	}

	return s.queue.Enqueue(constants.Decrease, req, source.ID)
}

func (s *CashbackService) GetCashbackByUserID(turonUserID int64) (*models.Cashback, error) {
	if err := s.validateTuronUserID(turonUserID); err != nil {
		return nil, err
	}
	return s.repo.GetCashbackByUserID(turonUserID)
}

func (s *CashbackService) GetCashbackHistoryByUserID(turonUserID int64, fromDate, toDate string, pagination *models.Pagination) ([]models.CashbackHistory, error) {
	if err := s.validateTuronUserID(turonUserID); err != nil {
		return nil, err
	}

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
