/*

 	baseUrl := "http://sandbox-cloudapi.imrworldwide.com/"
	bucket := "us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod"
	logParser := LogParser.New("EN660IH1TCIQT.2017-01-18-01.551a534d.gz", bucket, baseUrl)
	logParser.ReadObject()

	40 seconds
	170%

}
*/

package main

import (
	"context"
	"log"
	"me/logparser_lambda/logparser"
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
	verbose := os.Getenv("VERBOSE")

	if endPoint == "" {
		panic("Environment ENDPOINT needs to be defined")
	}

	for _, record := range s3Event.Records {
		//log.Println("RECORD", record)
		parser := logparser.NewParser(endPoint)
		if maxRequests != "" {
			nW, err := strconv.Atoi(maxRequests)
			if err != nil {
				panic("Don't understand MAX_REQUESTS = " + maxRequests)
			}
			nR, err := strconv.Atoi(maxWorkers)
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
