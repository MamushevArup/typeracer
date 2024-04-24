package main

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/handlers"
	"github.com/MamushevArup/typeracer/internal/lib/http/server"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/MamushevArup/typeracer/pkg/psql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	contextTimeout = 4 * time.Second
)

//	@title			Typeracer Game clone API
//	@version		2.0
//	@description	API for typeracer game clone. Typeracer popular game where users improve their typing skills in interactive format

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization

// @host		localhost:1001
// @schemes	http
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer func() {
		cancel()
	}()

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

	lg.Info("DB connection established")

	repo := repository.NewRepo(lg, db)

	svc := services.NewService(repo, cfg)

	handler := handlers.NewHandler(svc)

	lg.Info("Server started at port " + cfg.HttpServer.Port)

	// deactivate link under 1 hour usage go to the database every <duration>
	go svc.Link.Kill(time.NewTicker(10 * time.Second))

	srv, err := server.Http(handler, cfg)

	if err != nil {
		lg.Errorf("unable to create a connection %v", err)
		os.Exit(1)
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			lg.Errorf("unable to create a connection %v", err)
			os.Exit(1)
		}
	}()

	//Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctxShutDown, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		lg.Errorf("unable to shutdown server %v", err)
	}

	lg.Info("Server stopped")
}
