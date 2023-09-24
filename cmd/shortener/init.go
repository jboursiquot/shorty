package main

import (
	"os"

	"log/slog"

	"github.com/jboursiquot/shorty"
	"github.com/joeshaw/envdecode"
)

var (
	log *slog.Logger
	cfg shorty.Config
)

func init() {
	log = shorty.DefaultLogger()

	log.Info("Initializing Config...")
	cfg = shorty.Config{}
	if err := envdecode.Decode(&cfg); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
