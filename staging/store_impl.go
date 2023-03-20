package staging

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
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
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	staging := model.Staging{
		Table:  record.Table,
		Action: model.StagingActionCreate,
		Fields: record.Fields,
	}

	// Check if the record exist.
	pks, selects, query, args := record.Query()
	if err = conn.QueryRow(ctx, query, args...).Scan(selects...); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		return err
	}

	// Insert primary keys to fields if the record exist.
	if err == nil {
		staging.Action = model.StagingActionUpdate
		for i, pk := range pks {
			staging.Fields[pk] = sqlVarToAny(selects[i])
		}
	}

	// Insert the rest of the fields. If the field is a search pattern, then we search for the primary keys.
	for k, v := range record.Fields {
		switch v.(type) {
		case map[string]any:
			fieldJSON, err := json.Marshal(v)
			if err != nil {
				log.Println(err)
				return ErrorStagingBadInput
			}

			var r model.StagingCreateNestedSearch
			if err := json.Unmarshal(fieldJSON, &r); err != nil {
				log.Println(err)
				return ErrorStagingBadInput
			}

			if !r.Valid() {
				log.Println(r)
				return ErrorStagingBadInput
			}

			pks, selects, query, args := r.Query()
			if err := conn.QueryRow(ctx, query, args...).Scan(selects...); errors.Is(err, pgx.ErrNoRows) {
				return ErrorStagingFieldDepNotExist
			} else if err != nil {
				log.Println(err)
				return err
			}

			if len(pks) == 1 {
				staging.Fields[k] = sqlVarToAny(selects[0])
			} else {
				m := make(map[string]any)
				for i, pk := range pks {
					m[pk] = sqlVarToAny(selects[i])
				}
				staging.Fields[k] = m
			}
		case string:
		case float64:
		case bool:
		default:
			return ErrorStagingBadInput
		}
	}

	fieldsJSON, err := json.Marshal(staging.Fields)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := conn.Exec(ctx, `
		INSERT INTO staging_data (table_name, action, fields)
		VALUES ($1, $2, $3)
	`, staging.Table, staging.Action, fieldsJSON); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func sqlVarToAny(v any) any {
	switch v.(type) {
	case *sql.NullString:
		return v.(*sql.NullString).String
	case *sql.NullInt64:
		return v.(*sql.NullInt64).Int64
	case *sql.NullBool:
		return v.(*sql.NullBool).Bool
	case *sql.NullTime:
		return v.(*sql.NullTime).Time
	default:
		return nil
	}
}

func (s *impl) List(ctx context.Context, table model.StagingTable, offset, limit int) ([]*model.Staging, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	// Query from staging_data
	rows, err := conn.Query(ctx, `
		SELECT id, table_name, fields, action, created_at, updated_at
		FROM staging_data
		WHERE table_name = $1
		ORDER BY id DESC
		OFFSET $2 LIMIT $3
	`, table, offset, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	stagings := []*model.Staging{}
	for rows.Next() {
		var s model.Staging
		if err := rows.Scan(&s.Id, &s.Table, &s.Fields, &s.Action, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		stagings = append(stagings, &s)
	}
	if len(stagings) == 0 {
		return []*model.Staging{}, nil
	}

	// Generate query for existing records for compare
	conds := []string{}
	args := []any{}
	argsIdx := 1
	for _, s := range stagings {
		if s.Action != model.StagingActionUpdate {
			continue
		}

		ands := []string{}
		for _, pkName := range table.PkNames() {
			if _, ok := s.Fields[pkName]; !ok {
				return nil, ErrorStagingBadInput
			}

			ands = append(ands, fmt.Sprintf("%s = $%d", pkName, argsIdx))
			args = append(args, s.Fields[pkName])
			argsIdx++
		}

		conds = append(conds, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
	}

	if len(conds) == 0 {
		return stagings, nil
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE ", strings.Join(table.FieldNames(), ", "), table)
	query += strings.Join(conds, " OR ")

	olds := map[string]map[string]any{}
	rows, err = conn.Query(ctx, query, args...)
	for rows.Next() {
		fieldVars := table.FieldVars()
		if err := rows.Scan(fieldVars.Vars...); err != nil {
			log.Println(err)
			return nil, err
		}

		olds[fieldVars.KeyString()] = fieldVars.Map()
	}

	// Combine old and new records and return result
	for _, s := range stagings {
		if s.Action != model.StagingActionUpdate {
			continue
		}

		key := s.KeyString()
		m := map[string]any{}
		for _, fieldName := range table.FieldNames() {
			newVal, newOk := s.Fields[fieldName]
			oldVal, oldOk := olds[key][fieldName]

			compare := model.StagingFieldCompare{}
			if newOk {
				compare.New = newVal
			}
			if oldOk {
				compare.Old = oldVal
			}

			if fieldChanged(oldVal, newVal) {
				compare.Changed = true
			}

			m[fieldName] = compare
		}

		s.Fields = m
	}

	return stagings, nil
}

func fieldChanged(old, new any) bool {
	switch o := old.(type) {
	case int64:
		n, ok := new.(float64)
		if !ok {
			return true
		}
		return o != int64(n)
	case string:
		n, ok := new.(string)
		if !ok {
			return true
		}
		return o != n
	case bool:
		n, ok := new.(bool)
		if !ok {
			return true
		}
		return o != n
	}

	return true
}

// TODO implement the actual submit
func (s *impl) Submit(ctx context.Context, id int) error {
	// TODO Add transation
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

	return nil
}
