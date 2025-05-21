package core

import "cashback-serv/models"

type SourceFinderCreator interface {
	FindSourceOrCreate(turonUserID, cineramaUserID int64, hostIP string) (*models.Source, error)
}
