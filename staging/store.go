package staging

import (
	"context"
	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, record *model.StagingDataCreateRecord) error
	Submit(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*model.StagingData, error)
}
