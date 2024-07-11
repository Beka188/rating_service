package model

import "time"

type Rating struct {
	ID        int
	UserID    int
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
