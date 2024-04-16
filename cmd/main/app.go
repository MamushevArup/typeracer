package main

import (
	"context"
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

const (
	contextTimeout = 4 * time.Second
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't read .env %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("error due to: %v", err)
	}
	// postgres://username:password@localhost:5432/database_name

	lg := logger.New()

	db, err := psql.New(ctx, cfg)
	if err != nil {
		lg.Errorf("unable to create a connection %v", err)
		os.Exit(1)
	}

	defer db.Close()

	lg.Info("DB connection established")

	repo := repository.NewRepo(lg, db)

	svc := services.NewService(repo, cfg)

	handler := handlers.NewHandler(svc)

	lg.Info("Server started at port " + cfg.HttpServer.Port)

	// deactivate link under 1 hour usage go to the database every <duration>
	go svc.Link.Kill(time.NewTicker(10 * time.Second))

	go func() {
		if err = http.ListenAndServe(":"+cfg.HttpServer.Port, handler.InitRoutes()); err != nil {
			lg.Errorf("unable to create a connection %v", err)
			os.Exit(1)
		}
	}()

	select {}

}
