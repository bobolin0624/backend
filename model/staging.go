package model

import "time"

type StagingDataCreateRecord struct {
	Table    string              `json:"table"`
	SearchBy StagingDataSearchBy `json:"searchBy"`
	Fields   map[string]any      `json:"fields"`
}

type StagingDataSearchBy map[string]any

// TODO remove this
type StagingData struct {
	Id      int             `json:"id"`
	Records []StagingRecord `json:"records"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TODO remove this
type StagingRecord struct {
	Table  string                 `json:"table"`
	Record map[string]interface{} `json:"record"`
}
