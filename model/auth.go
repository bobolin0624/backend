package model

import (
	"errors"

	"google.golang.org/api/idtoken"
)

type AuthType int

const (
	AuthTypeGoogle AuthType = iota + 1
)

type AuthInfo struct {
	Type AuthType `json:"type"`

	Google *AuthInfoGoogle `json:"google"`
}

type AuthInfoGoogle struct {
	IdToken string `json:"idToken"`
}

type AuthResult struct {
	Type AuthType

	Google *AuthResultGoogle
}

func (r *AuthResult) ToUser() (*User, error) {
	switch r.Type {
	case AuthTypeGoogle:
		return r.Google.ToUser(), nil
	default:
		return nil, errors.New("unknown auth type")
	}
}

type AuthResultGoogle struct {
	Payload *idtoken.Payload
}

func (rg *AuthResultGoogle) ToUser() *User {
	return &User{
		Name:      rg.Payload.Claims["name"].(string),
		Email:     rg.Payload.Claims["email"].(string),
		AvatarURL: rg.Payload.Claims["picture"].(string),
		GoogleId:  rg.Payload.Subject,
	}
}
