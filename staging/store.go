package staging

import "context"

type Store interface {
	List(ctx context.Context, offset, limit int) ([]*StagingData, error)
}
