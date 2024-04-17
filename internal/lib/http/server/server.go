package server

import (
	"fmt"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/handlers"
	"net/http"
	"time"
)

func Http(handler handlers.Handler, cfg *config.Config) *http.Server {
	headerT, err := convertConfigToDuration(cfg.HttpServer.HeaderTimeout)
	if err != nil {
		return nil
	}

	idleT, err := convertConfigToDuration(cfg.HttpServer.IdleTimeout)
	if err != nil {
		return nil
	}

	readT, err := convertConfigToDuration(cfg.HttpServer.Timeout)
	if err != nil {
		return nil
	}

	return &http.Server{
		Addr:              ":" + cfg.HttpServer.Port,
		Handler:           handler.InitRoutes(),
		ReadHeaderTimeout: headerT,
		IdleTimeout:       idleT,
		ReadTimeout:       readT,
		MaxHeaderBytes:    1 << 20,
	}
}

func convertConfigToDuration(cfg string) (time.Duration, error) {
	res, err := time.ParseDuration(cfg)
	if err != nil {
		return 0, fmt.Errorf("unable to parse duration %w", err)
	}
	return res, nil
}
