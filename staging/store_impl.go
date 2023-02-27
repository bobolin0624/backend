package staging

import (
	"context"
	"errors"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (s *impl) Create(ctx context.Context, record *model.StagingDataCreateRecord) error {
	// Check if exist and return id and flag it update. If not flag it create.
	// Search for fields that needs searching for ids. If not found return failed.
	// Create fields and flags and insert into staging_data.
	return errors.New("TODO")
}

// TODO refactor
func (s *impl) List(ctx context.Context, offset, limit int) ([]*model.StagingData, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT id, records, created_at, updated_at
		FROM staging_data
		ORDER BY id DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		return nil, err
	}

	stagingData := []*model.StagingData{}
	for rows.Next() {
		var s model.StagingData
		if err := rows.Scan(&s.Id, &s.Records, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		stagingData = append(stagingData, &s)
	}

	return stagingData, nil
}

// TODO refactor
func (s *impl) Submit(ctx context.Context, id int) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err = conn.Exec(ctx, `
		DELETE FROM staging_data
		WHERE id = $1
	`, id); err != nil {
		return err
	}

	// TODO implement the actual submit

	return nil
}
