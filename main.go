package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/taiwan-voting-guide/backend/handler"
)

const userSessionName = "user_session"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading .env file failed")
	}

	store, err := initSession()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(sessions.Sessions(userSessionName, store))
	// public routes
	r.GET("/health", handler.HealthCheck)

	auth := r.Group("/auth")
	handler.MountAuthRoutes(auth)

	// authenticated routes
	workspace := r.Group("/workspace")
	handler.MountAdminRoutes(workspace)

	r.Run()
}

func initSession() (sessions.Store, error) {
	db, err := sql.Open("postgres", os.Getenv("PG_URL")+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	store, err := postgres.NewStore(db, []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return nil, err
	}

	return store, nil
}
