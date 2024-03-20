package main

import (
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/handlers"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/MamushevArup/typeracer/pkg/psql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

/*
	Type-racer main page
	1) Practice yourself
	2) Multiple connection url for the race
	3) User account
	4) Store models and set of models
	5) Top racers.
	6) Store users avatar (models like in reddit)
	7) Contribute a text

	Race mode
	1) Models and users maximum 5 user.
	2) Waiting when more 2 typist enter then start a race
	3) The random text which will be really random.
	4) Time for the race like 2 minutes.

	Practice yourself
	1) All the same except for the waiting.

	Sign-in
	1) Email, Password
	2) via Google

	Sign-up
	1) Email
	2) Password
	3) Username

	Profile page
	1) Average
	2) Best race
	3) Total races
	4) Avatar
	5) History

	Model's page
	1) Models and availability

	Avatar model's page
	1) Avatar models

	In the type race 1 stands for the practice yourself and 0 for racetrack multiple connection
*/

func main() {
	// Init .env reader
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't read .env %v", err)
	}

	// Init Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error with reading config %e", err)
	}
	// postgres://username:password@localhost:5432/database_name

	// Init logger
	lg := logger.NewLogger()
	// Init DBClient
	db := psql.NewDBConnector(cfg)

	defer db.Close()
	// Init repository
	repo := repository.NewRepo(lg, db)

	svc := services.NewService(repo)

	handler := handlers.NewHandler(svc)

	// deactivate link under 1 hour usage go to the database every <duration>
	go svc.Link.Kill(time.NewTicker(10 * time.Second))

	if err = http.ListenAndServe(":"+cfg.HttpServer.Port, handler.InitRoutes()); err != nil {
		lg.Errorf("unable to create a connection %v", err)
		os.Exit(1)
	}
	select {}

}
