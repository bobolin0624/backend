package main

import (
	"github.com/taiwan-voting-guide/backend/route"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", route.HealthCheck)
	r.Run()
}
