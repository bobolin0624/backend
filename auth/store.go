package auth

import (
	"context"
	"errors"

	"github.com/taiwan-voting-guide/backend/model"
)

var (
	ErrTypeInvalid          = errors.New("type invalid")
	ErrTokenAudienceInvalid = errors.New("token audience invalid")
)

type Store interface {
	Auth(ctx context.Context, info *model.AuthInfo) (*model.AuthResult, error)
}
