package model

import (
	"database/sql"
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

func (sc *StagingCreate) Valid() bool {
	return sc.Table.Valid() && sc.SearchBy.Valid() && sc.Fields.Valid()
}

func (sc *StagingCreate) Query() ([]string, []any, string, []any) {
	where := []string{}
	args := []any{}
	i := 1
	keys := []string{}
	for k := range sc.SearchBy {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		where = append(where, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, sc.SearchBy[k])
		i += 1
	}

	pks := sc.Table.Pks()
	selects := sc.Table.PkVars()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE "+strings.Join(where, " AND "), strings.Join(pks, ", "), sc.Table)
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

func (t StagingTable) Pks() []string {
	switch t {
	case StagingCreateTablePoliticians:
		return []string{"id"}
	case StagingCreateTableParties:
		return []string{"id"}
	case StagingCreateTableCandidates:
		return []string{"type", "term", "politician_id"}
	case StagingCreateTableLegislators:
		return []string{"politician_id", "term"}
	default:
		panic("unknown table")
	}
}

func (t StagingTable) PkIndex() []int {
	switch t {
	case StagingCreateTableParties:
		return []int{0}
	case StagingCreateTablePoliticians:
		return []int{0}
	case StagingCreateTableCandidates:
		return []int{0, 1, 2}
	case StagingCreateTableLegislators:
		return []int{0, 1}
	default:
		panic("unknown table")
	}
}

type FieldVars []any

func (t StagingTable) PkMap(vars FieldVars) map[string]any {
	pks := t.Pks()
	m := map[string]any{}
	for i, pk := range pks {
		switch v := vars[i].(type) {
		case *sql.NullInt64:
			if !v.Valid {
				continue
			}
			m[pk] = v.Int64
		case *sql.NullBool:
			if !v.Valid {
				continue
			}
			m[pk] = v.Bool
		case *sql.NullString:
			if !v.Valid {
				continue
			}
			m[pk] = v.String
		case *sql.NullTime:
			if !v.Valid {
				continue
			}
			m[pk] = v.Time
		}
	}
	return m

}

func (t StagingTable) PkKey(args FieldVars) string {
	pkMap := t.PkMap(args)
	keys := []string{}
	for _, pk := range t.Pks() {
		keys = append(keys, fmt.Sprintf("%v", pkMap[pk]))
	}
	return strings.Join(keys, "-")
}

func (t StagingTable) Map(vars FieldVars) map[string]any {
	fields := t.Fields()
	m := map[string]any{}
	for i, f := range fields {
		switch v := vars[i].(type) {
		case *sql.NullInt64:
			if !v.Valid {
				continue
			}
			m[f] = v.Int64
		case *sql.NullBool:
			if !v.Valid {
				continue
			}
			m[f] = v.Bool
		case *sql.NullString:
			if !v.Valid {
				continue
			}
			m[f] = v.String
		case *sql.NullTime:
			if !v.Valid {
				continue
			}
			m[f] = v.Time
		}
	}
	return m
}

func (t StagingTable) PkVars() FieldVars {
	switch t {
	case StagingCreateTableParties:
		var id int
		return FieldVars{&id}
	case StagingCreateTablePoliticians:
		var id int
		return FieldVars{&id}
	case StagingCreateTableCandidates:
		var t string
		var term, politicianId int
		return FieldVars{&t, &term, &politicianId}
	case StagingCreateTableLegislators:
		var politicianId, term int
		return FieldVars{&politicianId, &term}
	default:
		panic("unknown table")
	}
}

func (t StagingTable) Fields() []string {
	switch t {
	case StagingCreateTableParties:
		return []string{
			"id",
			"name",
			"chairman",
			"established_date",
			"filing_date",
			"main_office_address",
			"mailing_address",
			"phone_number",
			"status",
		}
	case StagingCreateTablePoliticians:
		return []string{
			"id",
			"name",
			"birthdate",
			"avatar_url",
			"sex",
			"current_party_id",
			"meta",
		}
	case StagingCreateTableCandidates:
		return []string{
			"type",
			"term",
			"politician_id",
			"number",
			"elected",
			"party_id",
			"area",
			"vice_president",
		}
	case StagingCreateTableLegislators:
		return []string{
			"politicians_id",
			"term",
			"party_id",
			"onboard_date",
			"resign_date",
			"resign_reason",
		}
	default:
		panic("unknown table")
	}
}

func (t StagingTable) FieldVars() FieldVars {
	switch t {
	case StagingCreateTableParties:
		var id sql.NullInt64
		var name, chairman, mainOfficeAddress, mailingAddress, phoneNumber, status sql.NullString
		var establishedDate, filingDate sql.NullTime
		return FieldVars{
			&id,
			&name,
			&chairman,
			&establishedDate,
			&filingDate,
			&mainOfficeAddress,
			&mailingAddress,
			&phoneNumber,
			&status,
		}
	case StagingCreateTablePoliticians:
		var id sql.NullInt64
		var name, avatarUrl, sex, currentPartyId sql.NullString
		var birthdate sql.NullTime
		return FieldVars{
			&id,
			&name,
			&birthdate,
			&avatarUrl,
			&sex,
			&currentPartyId,
		}
	case StagingCreateTableCandidates:
		var t, area sql.NullString
		var term, politicianId, number, partyId sql.NullInt64
		var elected, vicePresident sql.NullBool
		return FieldVars{
			&t,
			&term,
			&politicianId,
			&number,
			&elected,
			&partyId,
			&area,
			&vicePresident,
		}
	case StagingCreateTableLegislators:
		var politicianId, term, partyId sql.NullInt64
		var onboardDate, resignDate sql.NullTime
		var resignReason sql.NullString
		return FieldVars{
			&politicianId,
			&term,
			&partyId,
			&onboardDate,
			&resignDate,
			&resignReason,
		}
	default:
		panic("unknown table")
	}
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
	Action StagingAction  `json:"action"`
	Fields map[string]any `json:"fields"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Staging) PkKey() string {
	keys := []string{}
	for _, pk := range s.Table.Pks() {
		keys = append(keys, fmt.Sprintf("%v", s.Fields[pk]))
	}
	return strings.Join(keys, "-")
}

type StagingFieldCompare struct {
	Old any `json:"old"`
	New any `json:"new"`
}
