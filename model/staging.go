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

func (r *StagingCreate) Valid() bool {
	return r.Table.Valid() && r.SearchBy.Valid() && r.Fields.Valid()
}

func (r *StagingCreate) Query() ([]string, []any, string, []any) {
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

	pks := r.Table.Pks()
	selects := r.Table.PkVars()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE "+strings.Join(where, " AND "), strings.Join(pks, ", "), r.Table)
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

func (t StagingTable) PkVars() []any {
	switch t {
	case StagingCreateTableParties:
		var id int
		return []any{&id}
	case StagingCreateTablePoliticians:
		var id int
		return []any{&id}
	case StagingCreateTableCandidates:
		var t string
		var term, politicianId int
		return []any{&t, &term, &politicianId}
	case StagingCreateTableLegislators:
		var politicianId, term int
		return []any{&politicianId, &term}
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

func (t StagingTable) FieldVars() []any {
	switch t {
	case StagingCreateTableParties:
		var id int
		var name, chairman, mainOfficeAddress, mailingAddress, phoneNumber, status string
		var establishedDate, filingDate time.Time
		return []any{
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
		var id int
		var name, avatarUrl, sex, currentPartyId string
		var birthdate time.Time
		var meta []byte
		return []any{
			&id,
			&name,
			&birthdate,
			&avatarUrl,
			&sex,
			&currentPartyId,
			&meta,
		}
	case StagingCreateTableCandidates:
		var t, area string
		var term, politicianId, number, partyId int
		var elected, vicePresident bool
		return []any{
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
		var politicianId, term, partyId int
		var onboardDate, resignDate time.Time
		var resignReason string
		return []any{
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

type StagingFieldCompare struct {
	Old any `json:"old"`
	New any `json:"new"`
}
