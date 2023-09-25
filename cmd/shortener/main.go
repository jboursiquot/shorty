package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/jboursiquot/shorty"
)

func main() {
	urlFlag := flag.String("u", "", "URL to shorten")
	flag.Parse()

	db := shorty.NewDB(ddb, "shorty")
	s := shorty.NewShortener(db)

	if *urlFlag == "" {
		log.Error("Missing required flag: -u")
		return
	}

	ctx := context.Background()

	log.Info("Shortening URL...")
	key := s.Shorten(ctx, *urlFlag)
	url := fmt.Sprintf("%s/%s", cfg.BaseURL, key)
	log.Info("Shorten Results", "key", key, "url", url)

	log.Info("Resolving URL...")
	url = s.Resolve(ctx, key)
	log.Info("Resolve Results", "key", key, "url", url)
}
