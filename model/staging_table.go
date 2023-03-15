package model

import (
	"database/sql"
	"strconv"
	"strings"
)

var tables = map[StagingTable][]Field{
	StagingTableParties: {
		{"id", DataTypeNumber, true},
		{"name", DataTypeString, false},
		{"chairman", DataTypeString, false},
		{"established_date", DataTypeDate, false},
		{"filing_date", DataTypeDate, false},
		{"main_office_address", DataTypeString, false},
		{"mailing_address", DataTypeString, false},
		{"phone_number", DataTypeString, false},
		{"status", DataTypeString, false},
	},
	StagingTablePoliticians: {
		{"id", DataTypeNumber, true},
		{"name", DataTypeString, false},
		{"birthdate", DataTypeDate, false},
		{"avatar_url", DataTypeString, false},
		{"sex", DataTypeString, false},
		{"current_party_id", DataTypeNumber, false},
	},
	StagingTableCandidates: {
		{"type", DataTypeString, true},
		{"term", DataTypeNumber, true},
		{"politician_id", DataTypeNumber, true},
		{"number", DataTypeNumber, false},
		{"elected", DataTypeBool, false},
		{"party_id", DataTypeNumber, false},
		{"area", DataTypeString, false},
		{"vice_president", DataTypeBool, false},
	},
	StagingTableLegislators: {
		{"politician_id", DataTypeNumber, true},
		{"term", DataTypeNumber, true},
		{"party_id", DataTypeNumber, false},
		{"onboard_date", DataTypeDate, false},
		{"resign_date", DataTypeDate, false},
		{"resign_reason", DataTypeString, false},
	},
}

type StagingTable string

const (
	StagingTablePoliticians StagingTable = "politicians"
	StagingTableParties     StagingTable = "parties"
	StagingTableCandidates  StagingTable = "candidates"
	StagingTableLegislators StagingTable = "legislators"
)

func (t StagingTable) Valid() bool {
	_, ok := tables[t]
	return ok
}

func (t StagingTable) isField(str string) bool {
	for _, f := range t.Fields() {
		if f.Name == str {
			return true
		}
	}

	return false
}

func (t StagingTable) Fields() []Field {
	fields, _ := tables[t]

	fs := []Field{}
	for _, f := range fields {
		fs = append(fs, f)
	}
	return fs
}

func (t StagingTable) FieldNames() []string {
	fields := t.Fields()
	names := []string{}
	for _, f := range fields {
		names = append(names, f.Name)
	}
	return names
}

func (t StagingTable) FieldVars() FieldVars {
	fields, _ := tables[t]

	names := []string{}
	vars := []any{}
	for _, f := range fields {
		vars = append(vars, createVar(f.DataType))
		names = append(names, f.Name)
	}
	return FieldVars{Table: t, Names: names, Vars: vars}
}

func (t StagingTable) Pks() []Field {
	fields, _ := tables[t]

	pks := []Field{}
	for _, f := range fields {
		if f.Pk {
			pks = append(pks, f)
		}
	}
	return pks
}

func (t StagingTable) PkNames() []string {
	fields := t.Pks()
	names := []string{}
	for _, f := range fields {
		names = append(names, f.Name)
	}
	return names
}

func (t StagingTable) PkIndex() []int {
	fields, _ := tables[t]

	indexs := []int{}
	for i, f := range fields {
		if !f.Pk {
			continue
		}
		indexs = append(indexs, i)
	}
	return indexs
}

func (t StagingTable) PkVars() FieldVars {
	fields, _ := tables[t]

	vars := []any{}
	names := []string{}
	for _, f := range fields {
		if !f.Pk {
			continue
		}

		vars = append(vars, createVar(f.DataType))
		names = append(names, f.Name)
	}
	return FieldVars{Table: t, Names: names, Vars: vars}
}

func createVar(dt DataType) any {
	switch dt {
	case DataTypeString:
		return new(sql.NullString)
	case DataTypeNumber:
		return new(sql.NullInt64)
	case DataTypeBool:
		return new(sql.NullBool)
	case DataTypeDate:
		return new(sql.NullTime)
	default:
		panic("unknown data type")
	}
}

type DataType string

const (
	DataTypeString DataType = "string"
	DataTypeNumber DataType = "number"
	DataTypeBool   DataType = "boolean"
	DataTypeDate   DataType = "date"
)

type Field struct {
	Name     string
	DataType DataType
	Pk       bool
}

type FieldVars struct {
	Table StagingTable
	Names []string
	Vars  []any
}

func (fv FieldVars) Map() map[string]any {
	m := map[string]any{}
	for i, n := range fv.Names {
		switch v := fv.Vars[i].(type) {
		case *sql.NullInt64:
			if !v.Valid {
				continue
			}
			m[n] = v.Int64
		case *sql.NullBool:
			if !v.Valid {
				continue
			}
			m[n] = v.Bool
		case *sql.NullString:
			if !v.Valid {
				continue
			}
			m[n] = v.String
		case *sql.NullTime:
			if !v.Valid {
				continue
			}
			m[n] = v.Time.Format("2006-01-02")
		}
	}
	return m
}

func (fv FieldVars) KeyString() string {
	strs := []string{}
	for i := range fv.Table.PkIndex() {
		switch v := fv.Vars[i].(type) {
		case *sql.NullInt64:
			if !v.Valid {
				continue
			}
			strs = append(strs, strconv.FormatInt(v.Int64, 10))
		case *sql.NullBool:
			if !v.Valid {
				continue
			}
			strs = append(strs, strconv.FormatBool(v.Bool))
		case *sql.NullString:
			if !v.Valid {
				continue
			}
			strs = append(strs, v.String)
		case *sql.NullTime:
			if !v.Valid {
				continue
			}
			strs = append(strs, v.Time.Format("2006-01-02"))
		default:
		}
	}
	return strings.Join(strs, "-")
}
