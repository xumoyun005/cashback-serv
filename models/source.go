package models

import "time"

type Source struct {
	ID        int64      `json:"id"`
	HostIP    string     `json:"host_ip"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
