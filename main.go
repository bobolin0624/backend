package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/taiwan-voting-guide/backend/route"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading .env file failed")
	}

	r := gin.Default()
	r.GET("/health", route.HealthCheck)
	r.Run()
}
