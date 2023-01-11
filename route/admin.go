package route

import (
	"github.com/gin-gonic/gin"
	// "github.com/taiwan-voting-guide/backend/staging"
)

func MountAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/data", listStagingDataHandler)
	rg.POST("/data", submitStagingDataHandler)
}

func listStagingDataHandler(c *gin.Context) {
}

func submitStagingDataHandler(c *gin.Context) {
	c.JSON(501, gin.H{})
}
