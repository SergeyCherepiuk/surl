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

	sessionDeleter := redis.NewSessionDeleter()
	accountDeleter := postgres.NewAccountDeleter()

	originGetter := redis.NewOriginGetter(postgres.NewOriginGetter())
	urlGetter := postgres.NewUrlGetter()
	urlCreator := postgres.NewUrlCreator()
	urlUpdater := redis.NewUrlUpdater(postgres.NewUrlUpdater())
	urlDeleter := redis.NewUrlDeleter(postgres.NewUrlDeleter())

	e := http.Router{
		SessionChecker: sessionChecker,
		AccountGetter:  accountGetter,
		SessionCreator: sessionCreator,
		AccountCreator: accountCreator,
		SessionUpdater: sessionUpdater,
		AccountUpdater: accountUpdater,
		SessionDeleter: sessionDeleter,
		AccountDeleter: accountDeleter,
		OriginGetter:   originGetter,
		UrlGetter:      urlGetter,
		UrlCreator:     urlCreator,
		UrlUpdater:     urlUpdater,
		UrlDeleter:     urlDeleter,
	}.Build()
	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
