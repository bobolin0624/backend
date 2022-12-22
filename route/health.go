package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.Request.Context
	c.String(http.StatusOK, "healthy")
}
