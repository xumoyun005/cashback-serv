package models

import (
	"time"
)

type Cashback struct {
	ID             int64      `json:"id" db:"id" example:"1"`
	CashbackAmount float64    `json:"cashback_amount" db:"cashback_amount" example:"100.50"`
	TuronUserID    int64      `json:"turon_user_id" db:"turon_user_id" example:"123"`
	CineramaUserID int64      `json:"cinerama_user_id" db:"cinerama_user_id" example:"456"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at" example:"2024-03-20T10:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at" example:"2024-03-20T10:00:00Z"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at" example:"null"`
}

type CashbackHistory struct {
	ID             int64      `json:"id" db:"id" example:"1"`
	CashbackID     int64      `json:"cashback_id" db:"cashback_id" example:"1"`
	SourceID       int64      `json:"source_id" db:"source_id" example:"1"`
	SourceSlug     string     `json:"source_slug" db:"source_slug" example:"turon"`
	Type           string     `json:"type" db:"type" example:"turon"`
	CashbackAmount float64    `json:"cashback_amount" db:"cashback_amount" example:"50.25"`
	HostIP         string     `json:"host_ip" db:"host_ip" example:"192.168.1.1"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at" example:"2024-03-20T10:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at" example:"2024-03-20T10:00:00Z"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at" example:"null"`
}

type CashbackRequest struct {
	TuronUserID    int64   `json:"turon_user_id"`
	CineramaUserID int64   `json:"cinerama_user_id"`
	CashbackAmount float64 `json:"cashback_amount"`
	HostIP         string  `json:"host_ip"`
	Type           string  `json:"type"`
}
