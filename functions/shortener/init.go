package main

import (
	"log/slog"
	"os"

	"github.com/jboursiquot/shorty"
	"github.com/joeshaw/envdecode"
)

type config struct {
	shorty.Config
	Stage           string `env:"STAGE"`
	LocalServerPort string `env:"LOCAL_SERVER_PORT,default=8080"`
}

var (
	log *slog.Logger
	cfg config
)

func init() {
	log = shorty.DefaultLogger()

	log.Info("Initializing Config...")
	cfg = config{}
	if err := envdecode.Decode(&cfg); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
