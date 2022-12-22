package user

import (
	"errors"
	// "github.com/taiwan-voting-guide/backend/pg"

	"context"
)

var (
	errUserAlreadyExists = errors.New("user already exists")
	errUserNotFound      = errors.New("user not found")
)

type User struct {
	Id        string
	Name      string
	Email     string
	AvatarURL string

	GoogleId string
}

// Create a new user. If the user already exists, return errUserAlreadyExists.
func Create(ctx context.Context, user *User) error {
	return errors.New("TODO")
}

// Update an existing user. If the user does not exist, return errUserNotFound.
func Update(ctx context.Context, user *User) error {
	return errors.New("TODO")
}

// Get a user by id. If the user does not exist, return errUserNotFound.
func Get(ctx context.Context, id string) (*User, error) {
	return nil, errors.New("TODO")
}

// Soft delete a user by id. If the user does not exist, return errUserNotFound.
func Deactivate(ctx context.Context, id string) error {
	return errors.New("TODO")
}
