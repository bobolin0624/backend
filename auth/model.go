package auth

import (
	"google.golang.org/api/idtoken"
)

type Type int

const (
	TypeGoogle Type = iota + 1
)

type Info struct {
	Type Type `json:"type"`

	Google *InfoGoogle `json:"google"`
}

type InfoGoogle struct {
	IdToken string `json:"idToken"`
}

type Result struct {
	Type Type

	Google *ResultGoogle
}

type ResultGoogle struct {
	Payload *idtoken.Payload
}
