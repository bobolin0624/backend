package staging

import (
	"context"
	"errors"

	"github.com/taiwan-voting-guide/backend/model"
)

func NewStore() Store {
	return &impl{}
}

type impl struct{}

func (s *impl) List(ctx context.Context, offset, limit int) ([]*model.StagingData, error) {
	return nil, errors.New("TODO")
}
