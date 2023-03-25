package staging

import (
	"context"
	"errors"

	"github.com/taiwan-voting-guide/backend/model"
)

var (
	ErrorStagingBadInput              = errors.New("staging data bad input")
	ErrorStagingFieldDepNotExist      = errors.New("staging data field dependency not exists")
	ErrorStagingDuplicateSearchResult = errors.New("staging data duplicate search result")
	ErrorStagingInsertFailed          = errors.New("staging data insert failed")
)

type Store interface {
	Create(ctx context.Context, staging model.Staging) error
	Submit(ctx context.Context, id int, fields model.StagingFields) error
	List(ctx context.Context, table model.StagingTable, offset, limit int) ([]model.StagingResult, error)
	Delete(ctx context.Context, id int) error
}
