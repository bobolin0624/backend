package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/model"
)

func MustHavePermission(ctx context.Context, resource model.PermissionResource, actions []model.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		c.Next()
	}
}
