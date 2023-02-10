package user

import (
	"context"
	"errors"

	"github.com/taiwan-voting-guide/backend/model"
)

var (
	ErrUserExist       = errors.New("user exist")
	ErrUserNotFound    = errors.New("user not found")
	ErrAuthTypeInvalid = errors.New("auth type invalid")
)

type Store interface {
	CreateByAuthResult(ctx context.Context, authResult *model.AuthResult) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	GetByAuthResult(ctx context.Context, authResult *model.AuthResult) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
}
