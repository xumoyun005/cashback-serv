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

func (s *SourceService) FindSourceOrCreate(turonUserID, cineramaUserID int64, hostIP string) (*models.Source, error) {
	var slug string
	if turonUserID != 0 && cineramaUserID == 0 {
		slug = constants.SourceTuron
	} else if cineramaUserID != 0 && turonUserID == 0 {
		slug = constants.SourceCinerama
	} else {
		return nil, errors.New("cannot determine source slug: neither turon_user_id nor cinerama_user_id is provided, or both are provided")
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
