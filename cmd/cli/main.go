package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres"
	"github.com/SergeyCherepiuk/surl/pkg/database/redis"
	"github.com/SergeyCherepiuk/surl/pkg/http"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}
	postgres.MustConnect()
	redis.MustConnect()
}

func main() {
	// Services
	accountManagerService := postgres.NewAccountManagerService()
	sessionManagerService := redis.NewSessionManagerService()
	urlService := postgres.NewUrlService()

	e := http.NewRouter(accountManagerService, sessionManagerService, urlService)
	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
