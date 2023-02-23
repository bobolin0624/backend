package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("user_id") == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userId, exist := getSessionUserId(c); exist {
			c.Set("user_id", userId)
			return
		}

		c.Next()
	}
}

func getSessionUserId(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	userId := session.Get("user_id")

	if userId == nil {
		return "", false
	}

	return userId.(string), true
}
