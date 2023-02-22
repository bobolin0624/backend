package question

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

type impl struct{}

func New() Store {
	return &impl{}
}

func (im *impl) Create(ctx context.Context, q *model.PoliticianQuestionCreate) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, `
		INSERT INTO politician_questions (category, user_id, question, politician_id)
		VALUES ($1, $2, $3, $4)
	`, q.Category, q.UserId, q.Question, q.PoliticianId); err != nil {
		return err
	}

	return nil
}

func (im *impl) Search(ctx context.Context, politicianId int64, typ string) ([]*model.PoliticianQuestion, error) {
	return nil, nil
}
