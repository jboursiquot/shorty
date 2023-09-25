package main

import (
	"context"
	"log/slog"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	db  *shorty.DB
)

func init() {
	log = shorty.DefaultLogger()

	log.Info("Initializing Config...")
	cfg = config{}
	if err := envdecode.Decode(&cfg); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("Initializing DynamoDB Client...")
	c, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	ddb := dynamodb.NewFromConfig(c)
	db = shorty.NewDB(ddb, cfg.TableName)
}
