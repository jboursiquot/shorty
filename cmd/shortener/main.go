package main

import (
	"flag"
	"fmt"

	"github.com/jboursiquot/shorty"
)

func main() {
	urlFlag := flag.String("u", "", "URL to shorten")
	flag.Parse()

	s := shorty.NewShortener()

	if *urlFlag == "" {
		log.Error("Missing required flag: -u")
		return
	}

	key := s.Shorten(*urlFlag)
	url := fmt.Sprintf("%s/%s", cfg.BaseURL, key)
	log.Info("Results", "key", key, "url", url)
}
