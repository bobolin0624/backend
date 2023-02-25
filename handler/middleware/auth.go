package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	UserIdKey = "user_id"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userId, exist := getSessionUserId(c); exist {
			c.Set(UserIdKey, userId)
			return
		}

		c.Next()
	}
}

func MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(UserIdKey) == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func getSessionUserId(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	userId := session.Get(UserIdKey)

	if userId == nil {
		return "", false
	}

	return userId.(string), true
}
