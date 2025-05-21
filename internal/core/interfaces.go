package core

import "cashback-serv/models"

type SourceFinderCreator interface {
	FindSourceOrCreate(turonUserID int64, hostIP string) (*models.Source, error)
}
