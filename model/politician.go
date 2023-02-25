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
	AvatarUrl string `json:"avatarUrl"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}
