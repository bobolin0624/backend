package model

import "time"

type Party struct {
	Id                int
	Name              string
	Chairman          string
	EstablishedDate   time.Time
	FilingDate        time.Time
	MainOfficeAddress string
	MailingAddress    string
	PhoneNumber       string
	Status            string

	CreatedAt time.Time
	UpdatedAt time.Time
}
