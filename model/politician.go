package model

import (
	"time"
)

type Politician struct {
	Id        int64
	Name      string
	Birthdate string
	AvatarUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Politician) Repr() *PoliticianRepr {
	return &PoliticianRepr{
		Name:      p.Name,
		Birthdate: p.Birthdate,
		AvatarUrl: p.AvatarUrl,

		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
}

type PoliticianRepr struct {
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	AvatarUrl string `json:"avatar_url"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
