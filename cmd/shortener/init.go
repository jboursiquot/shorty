package main

import (
	"context"
	"os"

	"log/slog"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jboursiquot/shorty"
	"github.com/joeshaw/envdecode"
)

var (
	log *slog.Logger
	cfg shorty.Config
	ddb shorty.DDBClient
)

func init() {
	log = shorty.DefaultLogger()

	log.Info("Initializing Config...")
	cfg = shorty.Config{}
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
	ddb = dynamodb.NewFromConfig(c)
}
