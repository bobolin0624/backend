package staging

import (
	"context"
	"errors"

	"github.com/taiwan-voting-guide/backend/model"
)

var (
	ErrorStagingNoChange         = errors.New("staging data no change")
	ErrorStagingBadInput         = errors.New("staging data bad input")
	ErrorStagingFieldDepNotExist = errors.New("staging data field dependency not exists")
)

type Store interface {
	Create(ctx context.Context, record *model.StagingCreate) error
	Submit(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*model.Staging, error)
}
