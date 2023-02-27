package question

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, q *model.PoliticianQuestionCreate) error
	List(ctx context.Context, politicianId int, offset, limit int) ([]*model.PoliticianQuestion, error)
	Search(ctx context.Context, politicianId int, typ string) ([]*model.PoliticianQuestion, error)
}
