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
	sessionChecker := redis.NewSessionChecker()
	accountGetter := postgres.NewAccountGetter()

	accountCreator := postgres.NewAccountCreator()
	sessionCreator := redis.NewSessionCreator()

	accountUpdater := postgres.NewAccountUpdater()
	sessionUpdater := redis.NewAccountUpdater(accountUpdater)

	accountDeleter := redis.NewAccountDeleter(postgres.NewAccountDeleter())

	urlService := postgres.NewUrlService()

	e := http.Router{
		SessionChecker: sessionChecker,
		AccountGetter:  accountGetter,
		AccountCreator: accountCreator,
		SessionCreator: sessionCreator,
		AccountUpdater: accountUpdater,
		SessionUpdater: sessionUpdater,
		AccountDeleter: accountDeleter,
		UrlService:     urlService,
	}.Build()
	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
