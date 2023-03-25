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

func (s *impl) Create(ctx context.Context, staging model.Staging) error {
	if ok, err := staging.Valid(); !ok {
		return err
	}

	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	searchByJSON, err := json.Marshal(staging.SearchBy)
	if err != nil {
		return err
	}

	fields := model.StagingFields{}

	// Insert the rest of the fields. If the field is a search pattern, then we search for the primary keys.
	for k, v := range staging.Fields {
		switch v.(type) {
		case map[string]any:
			fieldJSON, err := json.Marshal(v)
			if err != nil {
				return err
			}

			var ns model.StagingNestedSearch
			if err := json.Unmarshal(fieldJSON, &ns); err != nil {
				return err
			}

			if ok, err := ns.Valid(); !ok {
				return err
			}

			pks, selects, query, args := ns.Query()
			if err := conn.QueryRow(ctx, query, args...).Scan(selects...); errors.Is(err, pgx.ErrNoRows) {
				return ErrorStagingFieldDepNotExist
			} else if err != nil {
				return err
			}

			if len(pks) == 1 {
				fields[k] = sqlVarToAny(selects[0])
				continue
			}

			m := model.StagingFields{}
			for i, pk := range pks {
				m[pk] = sqlVarToAny(selects[i])
			}
			fields[k] = m

		case string:
			fields[k] = v
		case float64:
			fields[k] = v
		case bool:
			fields[k] = v
		default:
			return ErrorStagingBadInput
		}
	}

	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	if _, err := conn.Exec(ctx, `
		INSERT INTO staging_data (table_name, search_by, fields)
		VALUES ($1, $2, $3)
		ON CONFLICT (table_name, search_by, fields)
		DO UPDATE SET updated_at = NOW()
	`, staging.Table, searchByJSON, fieldsJSON); err != nil {
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

func (s *impl) List(ctx context.Context, table model.StagingTable, offset, limit int) ([]model.StagingResult, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	// Query from staging_data
	rows, err := conn.Query(ctx, `
		SELECT id, table_name, search_by, fields, created_at 
		FROM staging_data
		WHERE table_name = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`, table, offset, limit)
	if err != nil {
		return nil, err
	}

	stagings := []*model.Staging{}
	for rows.Next() {
		var s model.Staging
		if err := rows.Scan(&s.Id, &s.Table, &s.SearchBy, &s.Fields, &s.CreatedAt); err != nil {
			return nil, err
		}

		stagings = append(stagings, &s)
	}

	if len(stagings) == 0 {
		return []model.StagingResult{}, nil
	}

	// dedup searchBys
	searchBys := []model.StagingFields{}
	for _, staging := range stagings {
		for _, searchBy := range searchBys {
			if searchBy.Equal(staging.SearchBy) {
				break
			}
		}
		searchBys = append(searchBys, staging.SearchBy)
	}

	// create select condition by searchBys
	conds := []string{}
	args := []any{}
	argsIdx := 1
	for _, searchBy := range searchBys {
		ands := []string{}
		for field, value := range searchBy {
			ands = append(ands, fmt.Sprintf("%s = $%d", field, argsIdx))
			args = append(args, value)
			argsIdx++
		}
		conds = append(conds, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
	}

	if len(conds) == 0 {
		return []model.StagingResult{}, nil
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE ", strings.Join(table.FieldNames(), ", "), table)
	query += strings.Join(conds, " OR ")

	oldFieldsList := []model.StagingFields{}
	rows, err = conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		fieldVars := table.FieldVars()
		if err := rows.Scan(fieldVars.Vars...); err != nil {
			log.Println(err)
			return nil, err
		}

		oldFieldsList = append(oldFieldsList, fieldVars.Map())
	}

	// Combine old and new records and return result
	results := []model.StagingResult{}
	for _, s := range stagings {
		foundFields := []model.StagingFields{}
		for _, oldFields := range oldFieldsList {
			if s.SearchBy.ExistIn(oldFields) {
				foundFields = append(foundFields, oldFields)
			}
		}

		switch len(foundFields) {
		case 0:
			resultFields := []model.StagingResultField{}
			for _, fn := range table.FieldNames() {
				v := s.Fields[fn]
				resultFields = append(resultFields, model.StagingResultField{
					Type:  model.StagingResultFieldTypeValue,
					Name:  fn,
					Value: v,
				})
			}

			results = append(results, model.StagingResult{
				Id:     s.Id,
				Fields: resultFields,
				Status: model.StagingResultStatusCreate,
			})
		case 1:
			// Inject primary key from search result
			oldFields := foundFields[0]
			for _, pk := range table.PkNames() {
				s.Fields[pk] = oldFields[pk]
			}

			changed := false
			resultFields := []model.StagingResultField{}
			for _, fieldName := range table.FieldNames() {
				newVal, newOk := s.Fields[fieldName]
				oldVal, oldOk := oldFields[fieldName]

				compare := model.StagingFieldCompare{}
				if newOk {
					compare.New = newVal
				}
				if oldOk {
					compare.Old = oldVal
				}

				if fieldChanged(oldVal, newVal) {
					compare.Changed = true
					changed = true
				}

				resultFields = append(resultFields, model.StagingResultField{
					Type:  model.StagingResultFieldTypeCompare,
					Name:  fieldName,
					Value: compare,
				})
			}

			status := model.StagingResultStatusUpdate
			if !changed {
				status = model.StagingResultStatusDuplicate
			}

			results = append(results, model.StagingResult{
				Id:     s.Id,
				Fields: resultFields,
				Status: status,
			})
		default:
			resultFields := []model.StagingResultField{}
			for _, fn := range table.FieldNames() {
				if v, ok := s.Fields[fn]; ok {
					resultFields = append(resultFields, model.StagingResultField{
						Type:  model.StagingResultFieldTypeValue,
						Name:  fn,
						Value: v,
					})
				}
			}
			results = append(results, model.StagingResult{
				Id:     s.Id,
				Fields: resultFields,
				Status: model.StagingResultStatusConflict,
			})
		}

	}

	return results, nil
}

func fieldChanged(old, new any) bool {
	if new == nil {
		return false
	}

	switch o := old.(type) {
	case int64:
		switch n := new.(type) {
		case int64:
			return o != n
		case float64:
			return float64(o) != n
		}
	case float64:
		switch n := new.(type) {
		case int64:
			return o != float64(n)
		case float64:
			return o != n
		}
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

func (s *impl) Submit(ctx context.Context, id int, fields model.StagingFields) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	var table model.StagingTable
	searchBy := model.StagingFields{}
	err = tx.QueryRow(ctx, `
		SELECT table_name, search_by 
		FROM staging_data
		WHERE id = $1
		FOR UPDATE
	`, id).Scan(&table, &searchBy)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	ands := []string{}
	args := []any{}
	argsIdx := 1
	for field, value := range searchBy {
		ands = append(ands, fmt.Sprintf("%s = $%d", field, argsIdx))
		args = append(args, value)
		argsIdx++
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE ", strings.Join(table.FieldNames(), ", "), table)
	query += strings.Join(ands, " AND ")
	query += " FOR UPDATE"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	oldFieldsList := []model.StagingFields{}
	for rows.Next() {
		fieldVars := table.FieldVars()
		if err := rows.Scan(fieldVars.Vars...); err != nil {
			log.Println(err)
			return err
		}

		oldFieldsList = append(oldFieldsList, fieldVars.Map())
	}

	insertFields := []string{}
	insertValues := []string{}
	valueIdx := 1
	insertVars := []any{}
	for _, field := range table.FieldNames() {
		if _, ok := fields[field]; !ok {
			continue
		}

		insertFields = append(insertFields, field)
		insertValues = append(insertValues, fmt.Sprintf("$%d", valueIdx))
		insertVars = append(insertVars, fields[field])
		valueIdx++
	}

	switch len(oldFieldsList) {
	case 0:
		inserts := strings.Join(insertFields, ", ")
		values := strings.Join(insertValues, ", ")

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, inserts, values)
		tag, err := tx.Exec(ctx, query, insertVars...)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		if rowsAffected := tag.RowsAffected(); rowsAffected != 1 {
			tx.Rollback(ctx)
			return ErrorStagingInsertFailed
		}
	case 1:
		oldFields := oldFieldsList[0]
		sets := []string{}
		for i, field := range insertFields {
			sets = append(sets, fmt.Sprintf("%s = $%d", field, i+1))
		}

		where := []string{}
		for _, field := range table.PkNames() {
			where = append(where, fmt.Sprintf("%s = $%d", field, valueIdx))
			insertVars = append(insertVars, oldFields[field])
			valueIdx++
		}

		query := fmt.Sprintf("UPDATE %s SET %s WHERE ", table, strings.Join(sets, ", "))
		query += strings.Join(where, " AND ")
		log.Println(query)

	default:
		tx.Rollback(ctx)
		return ErrorStagingDuplicateSearchResult
	}

	if _, err = tx.Exec(ctx, `
		DELETE FROM staging_data
		WHERE id = $1
	`, id); err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}

func (s *impl) Delete(ctx context.Context, id int) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `
		DELETE FROM staging_data
		WHERE id = $1
	`, id)

	return err
}
