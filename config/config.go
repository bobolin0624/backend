package config

import (
	"os"
)

func GetFrontendHost() string {
	return os.Getenv("FRONTEND_HOST")
}

func GetFrontendEndpoint() string {
	return "http://" + os.Getenv("FRONTEND_HOST") + ":" + os.Getenv("FRONTEND_PORT")
}
