//
//
// Lambda Wrapper for logParser
//

package main

import (
	"context"
	"log"
	"logparser_lambda/logparser"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var invokeCount = 0

func init() {
}

func LambdaHandler(ctx context.Context, s3Event events.S3Event) (int, error) {
	invokeCount = invokeCount + 1
	maxRequests := os.Getenv("MAX_REQUESTS")
	maxWorkers := os.Getenv("MAX_WORKERS")
	endPoint := os.Getenv("ENDPOINT")

	if endPoint == "" {
		panic("Environment ENDPOINT needs to be defined")
	}

	for _, record := range s3Event.Records {
		//log.Println("RECORD", record)
		parser := logparser.NewParser(endPoint)
		if maxRequests != "" {
			nR, err := strconv.Atoi(maxRequests)
			if err != nil {
				panic("Don't understand MAX_REQUESTS = " + maxRequests)
			}
			nW, err := strconv.Atoi(maxWorkers)
			if err != nil {
				panic("Don't understand MAX_WORKERS = " + maxWorkers)
			}
			parser.WithMaxRequest(nR)
			parser.WithMaxWorkers(nW)
		}
		parser.ParseS3Object(record.S3)

		log.Println("LOGPARSER STATS: ", record.S3.Object.Key, "  = ", parser.GetStats())
	}

	return invokeCount, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
