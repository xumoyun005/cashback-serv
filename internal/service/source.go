package service

import (
	constants "cashback-serv/const"
	"cashback-serv/models"
	"errors"
	"fmt"
)

type SourceRepository interface {
	CreateSource(source *models.Source) error
	GetSourceBySlug(slug string) (*models.Source, error)
}

type SourceService struct {
	repo SourceRepository
}

func NewSourceService(repo SourceRepository) *SourceService {
	return &SourceService{repo: repo}
}

func (s *SourceService) FindSourceOrCreate(turonUserID int64, hostIP string) (*models.Source, error) {
	var slug string
	if turonUserID != 0 {
		slug = constants.SourceTuron
	} else {
		return nil, errors.New("cannot determine source slug: turon_user_id must be provided")
	}

	source, err := s.repo.GetSourceBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get source by slug: %w", err)
	}
	if source != nil {
		return source, nil
	}

	newSource := &models.Source{
		HostIP: hostIP,
		Slug:   slug,
	}

	err = s.repo.CreateSource(newSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	return newSource, nil
}
