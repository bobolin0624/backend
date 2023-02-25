package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/handler/middleware"
)

func HealthCheck(c *gin.Context) {
	userId := c.GetString(middleware.UserIdKey)

	loggedInMsg := "You are not logged in."
	if userId != "" {
		loggedInMsg = fmt.Sprintf("You are logged in as %s.", userId)
	}

	c.String(http.StatusOK, fmt.Sprintf("Healthy! %s", loggedInMsg))
}
