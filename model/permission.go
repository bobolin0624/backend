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

func PermissionPolitician(id int) Resource {
	return Resource("politician:" + strconv.FormatInt(int64(id), 10))
}

type Action string

const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
)

type PermissionResource string
