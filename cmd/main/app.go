package main

import (
	"fmt"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/repository/postgres"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/MamushevArup/typeracer/pkg/psql"
	"log"
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
*/

func main() {
	// Init Config
	cfg, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("error with reading config %e", err)
	}
	// postgres://username:password@localhost:5432/database_name
	pg := cfg.Postgres
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pg.User, pg.Password, pg.Host, pg.Port, pg.Database)
	// Init logger
	lg := logger.NewLogger()
	fmt.Println(dbUrl)
	// Init DBClient
	dbConn := psql.NewDBConnector(lg, dbUrl)
	// Init repository  TODO change the logic here
	dbConn.DBConnector()
	postgres.NewRepository(lg, dbConn)
}
