package question

import (
	"context"
	"errors"

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
	return nil, errors.New("not implemented")
}

func (im *impl) List(ctx context.Context, politicianId int64, offset, limit int) ([]*model.PoliticianQuestion, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT pq.id, pq.category, u.name , pq.question, pq.asked_at, pq.politician_id, pq.reply, pq.replied_at, pq.likes 
		FROM politician_questions pq
		INNER JOIN users u ON user_id = u.id
		WHERE politician_id = $1
		AND hidden = false
		ORDER BY asked_at DESC
		OFFSET $2 LIMIT $3
	`, politicianId, offset, limit)

	if err != nil {
		return nil, err
	}
		
	questions := []*model.PoliticianQuestion{}
	for rows.Next() {
		q := model.PoliticianQuestion{}
		if err := rows.Scan(&q.Id, &q.Category, &q.UserName, &q.Question, &q.AskedAt, &q.PoliticianId, &q.Reply, &q.RepliedAt, &q.Likes); err != nil {
			return nil, err
		}
		questions = append(questions, &q)
	}

	return questions, nil
}
