// Package main is the BlobTube server entry point.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sixfeetup/blobtube/internal/config"
	"github.com/sixfeetup/blobtube/internal/server"
)

func main() {
	cfg := config.FromEnv()

	zerolog.TimeFieldFormat = time.RFC3339Nano
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	log.Logger = zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Run(ctx, cfg, server.WithSignals(os.Interrupt, syscall.SIGTERM)); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		log.Fatal().Err(err).Msg("server exited")
	}
}
