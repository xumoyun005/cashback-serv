package queue

import (
	constants "cashback-serv/const"
	"cashback-serv/internal/interfaces"
	"cashback-serv/models"
	"errors"
	"sync"
)

type CashbackRepository interface {
	GetCashbackByUserID(userID int64, fromDate, toDate string) (*models.Cashback, error)
	CreateCashback(cashback *models.Cashback) error
	UpdateCashbackAmount(id int64, amount float64) error
	CreateCashbackHistory(history *models.CashbackHistory) error
}

type SourceFinderCreator interface {
	FindSourceOrCreate(turonUserID, cineramaUserID int64, hostIP string) (*models.Source, error)
}

type QueueRequest struct {
	*models.CashbackRequest
	SourceID int64
}

type CashbackOperation struct {
	Type     string
	Request  *QueueRequest
	Response chan error
}

type CashbackQueue struct {
	operations    chan *CashbackOperation
	userLocks     map[int64]*sync.Mutex
	mu            sync.RWMutex
	repo          CashbackRepository
	sourceService core.SourceFinderCreator
}

func NewCashbackQueue(repo CashbackRepository, sourceService core.SourceFinderCreator) *CashbackQueue {
	queue := &CashbackQueue{
		operations:    make(chan *CashbackOperation, 1000),
		userLocks:     make(map[int64]*sync.Mutex),
		repo:          repo,
		sourceService: sourceService,
	}

	go queue.process()

	return queue
}

func (q *CashbackQueue) process() {
	for op := range q.operations {
		q.mu.Lock()
		lock, exists := q.userLocks[op.Request.TuronUserID]
		if !exists {
			lock = &sync.Mutex{}
			q.userLocks[op.Request.TuronUserID] = lock
		}
		q.mu.Unlock()

		lock.Lock()
		var err error

		switch op.Type {
		case constants.Increase:
			err = q.handleIncrease(op.Request)
		case constants.Decrease:
			err = q.handleDecrease(op.Request)
		default:
			err = errors.New("unknown operation type")
		}

		lock.Unlock()
		op.Response <- err
	}
}

func (q *CashbackQueue) handleIncrease(req *QueueRequest) error {
	cashback, err := q.repo.GetCashbackByUserID(req.TuronUserID, "", "")
	if err != nil {
		return err
	}

	if cashback == nil {
		cashback = &models.Cashback{
			CashbackAmount: req.CashbackAmount,
			TuronUserID:    req.TuronUserID,
			CineramaUserID: req.CineramaUserID,
		}
		if err := q.repo.CreateCashback(cashback); err != nil {
			return err
		}
	} else {
		newAmount := cashback.CashbackAmount + req.CashbackAmount
		if err := q.repo.UpdateCashbackAmount(cashback.ID, newAmount); err != nil {
			return err
		}
		cashback.CashbackAmount = newAmount
	}

	history := &models.CashbackHistory{
		CashbackID:     cashback.ID,
		SourceID:       req.SourceID,
		CashbackAmount: req.CashbackAmount,
		HostIP:         req.HostIP,
		Type:           req.Type,
	}
	return q.repo.CreateCashbackHistory(history)
}

func (q *CashbackQueue) handleDecrease(req *QueueRequest) error {
	cashback, err := q.repo.GetCashbackByUserID(req.TuronUserID, "", "")
	if err != nil {
		return err
	}

	if cashback == nil {
		return errors.New("no cashback found for user")
	}

	if cashback.CashbackAmount < req.CashbackAmount {
		return errors.New("insufficient cashback amount")
	}

	newAmount := cashback.CashbackAmount - req.CashbackAmount
	if err := q.repo.UpdateCashbackAmount(cashback.ID, newAmount); err != nil {
		return err
	}

	history := &models.CashbackHistory{
		CashbackID:     cashback.ID,
		SourceID:       req.SourceID,
		CashbackAmount: req.CashbackAmount,
		HostIP:         req.HostIP,
		Type:           req.Type,
	}
	return q.repo.CreateCashbackHistory(history)
}

func (q *CashbackQueue) Enqueue(opType string, req *models.CashbackRequest, sourceID int64) error {
	queueReq := &QueueRequest{
		CashbackRequest: req,
		SourceID:        sourceID,
	}

	op := &CashbackOperation{
		Type:     opType,
		Request:  queueReq,
		Response: make(chan error, 1),
	}

	q.operations <- op

	return <-op.Response
}
