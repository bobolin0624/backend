package question

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, q *model.PoliticianQuestionCreate) error
	Search(ctx context.Context, politicianId int64, typ string) ([]*model.PoliticianQuestion, error)
}
