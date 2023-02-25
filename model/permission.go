package model

import (
	"strconv"
	"time"
)

type Permission struct {
	UserId   string
	Resource Resource
	Action   Action

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Resource string

const (
	ResouceAdmin    Resource = "admin"
	ResourceStaging Resource = "staging"
)

func GetPoliticianResource(id int64) Resource {
	return Resource("politician:" + strconv.FormatInt(id, 10))
}

type Action string

const (
	ActionAll    Action = "all"
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
)

type PermissionResource string
