package models

import "github.com/google/uuid"

type Event struct {
	ID       uuid.UUID
	SenderID int64
	Time     int64
	Name     string
}
