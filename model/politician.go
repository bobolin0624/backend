package model

import (
	"time"
)

type Politician struct {
	Id        int64
	Name      string
	Birthdate string
	EnName    string
	AvatarUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
}
