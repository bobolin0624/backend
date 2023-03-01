package model

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type StagingCreate struct {
	Table    StagingTable          `json:"table"`
	SearchBy StagingCreateSearchBy `json:"searchBy"`
	Fields   StagingCreateFields   `json:"fields"`
}

var (
	stagingPrimaryKey = map[StagingTable][]string{
		StagingCreateTableParties:     {"id"},
		StagingCreateTablePoliticians: {"id"},
		StagingCreateTableCandidates:  {},
		StagingCreateTableLegislators: {"id"},
	}
)

func (r *StagingCreate) Valid() bool {
	return r.Table.Valid() && r.SearchBy.Valid() && r.Fields.Valid()
}

func (r *StagingCreate) CreateQuery() ([]string, []any, string, []any) {
	where := []string{}
	args := []any{}
	i := 1
	keys := []string{}
	for k := range r.SearchBy {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		where = append(where, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, r.SearchBy[k])
		i += 1
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE "+strings.Join(where, " AND "), r.Table)
	pks, selects := r.Table.CreatePrimaryKeyVars()
	return pks, selects, query, args
}

type StagingTable string

const (
	StagingCreateTablePoliticians StagingTable = "politicians"
	StagingCreateTableParties     StagingTable = "parties"
	StagingCreateTableCandidates  StagingTable = "candidates"
	StagingCreateTableLegislators StagingTable = "legislators"
)

func (t StagingTable) Valid() bool {
	return t == StagingCreateTablePoliticians ||
		t == StagingCreateTableParties ||
		t == StagingCreateTableCandidates ||
		t == StagingCreateTableLegislators
}

func (t StagingTable) CreatePrimaryKeyVars() ([]string, []any) {
	switch t {
	case StagingCreateTablePoliticians:
		var id int
		return []string{"id"}, []any{&id}
	case StagingCreateTableParties:
		var id int
		return []string{"id"}, []any{&id}
	case StagingCreateTableCandidates:
		var t string
		var term, politicianId int
		return []string{"type", "term", "politician_id"}, []any{&t, &term, &politicianId}
	case StagingCreateTableLegislators:
		var id int
		return []string{"id"}, []any{&id}
	}
	return nil, nil
}

type StagingCreateSearchBy map[string]any

func (s StagingCreateSearchBy) Valid() bool {
	for _, v := range s {
		switch v.(type) {
		case int:
		case string:
		default:
			return false
		}
	}

	return true
}

type StagingCreateFields map[string]any

func (f StagingCreateFields) Valid() bool {
	if len(f) == 0 {
		return false
	}
	for _, v := range f {
		switch v.(type) {
		case map[string]any:
			return StagingCreateFields(v.(map[string]any)).Valid()
		case int:
			return true
		case string:
			return true
		}
	}

	return false
}

type StagingAction string

const (
	StagingActionCreate StagingAction = "create"
	StagingActionUpdate StagingAction = "update"
)

type Staging struct {
	Id     int            `json:"id"`
	Table  StagingTable   `json:"table"`
	Fields map[string]any `json:"fields"`
	Action StagingAction  `json:"action"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
