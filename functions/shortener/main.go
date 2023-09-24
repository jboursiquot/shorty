package main

import (
	"context"
	"flag"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

func main() {
	runAsLambda := flag.Bool("run-as-lambda", true, "Run as a lambda function")
	flag.Parse()

	router, err := newRouter(&cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	if *runAsLambda {
		ginLambda := ginadapter.New(router.Engine)
		handler := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return ginLambda.ProxyWithContext(ctx, req)
		}
		lambda.Start(handler)
	} else {
		if err := router.Run(":" + cfg.LocalServerPort); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}
}
