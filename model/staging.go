package model

import "time"

type StagingData struct {
	Id      int             `json:"id"`
	Records []StagingRecord `json:"records"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type StagingRecord struct {
	Table  string                 `json:"table"`
	Record map[string]interface{} `json:"record"`
}
