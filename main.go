package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/taiwan-voting-guide/backend/handler"
	"github.com/taiwan-voting-guide/backend/handler/middleware"
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
	// TODO set a proper CORS policy
	r.Use(cors.Default())
	r.Use(sessions.Sessions(userSessionName, store))
	r.Use(middleware.Auth())
	r.GET("/health", handler.HealthCheck)

	handler.MountAuthRoutes(r.Group("/auth"))
	handler.MountUserRoutes(r.Group("/user"))
	handler.MountWorkspaceRoutes(r.Group("/workspace"))
	handler.MountPolitician(r.Group("/politician"))

	r.Run()
}

func initSession() (sessions.Store, error) {
	db, err := sql.Open("postgres", os.Getenv("PG_URL"))
	if err != nil {
		return nil, err
	}

	store, err := postgres.NewStore(db, []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return nil, err
	}

	return store, nil
}
