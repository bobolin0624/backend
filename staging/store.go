package staging

import (
	"context"
	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	List(ctx context.Context, offset, limit int) ([]*model.StagingData, error)
	Submit(ctx context.Context, id int64) error
}
