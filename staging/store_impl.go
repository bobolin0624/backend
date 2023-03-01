package staging

import (
	"context"
	"errors"
	"log"

	// "github.com/jackc/pgx/v5"

	"github.com/jackc/pgx"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (s *impl) Create(ctx context.Context, record *model.StagingCreate) error {
	if !record.Valid() {
		return ErrorStagingBadInput
	}
	// Check if exist and return id and flag it update. If not flag it create.
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	staging := model.Staging{}
	_, selects, query, args := record.CreateQuery()
	if err := conn.QueryRow(ctx, query, args...).Scan(selects...); err == pgx.ErrNoRows {
		staging.Action = model.StagingActionCreate
	} else if err != nil {
		log.Println(err)
		return err
	}

	// Search for fields that needs searching for ids. If not found return failed.
	// Create fields and flags and insert into staging_data if some fields changes
	return errors.New("TODO")
}

// TODO add diff
func (s *impl) List(ctx context.Context, offset, limit int) ([]*model.Staging, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT id, table, fields, action, created_at, updated_at
		FROM staging_data
		ORDER BY id DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		return nil, err
	}

	staging := []*model.Staging{}
	for rows.Next() {
		var s model.Staging
		if err := rows.Scan(&s.Id, &s.Table, &s.Fields, &s.Action, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		staging = append(staging, &s)
	}

	return staging, nil
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
