package main

import (
	"encoding/json"
	"log"
	"logparser_lambda/logparser"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestIntegration(*testing.T) {
	log.Println("TestIntegration")
	var parser = logparser.NewParser(
		"http://sandbox-cloudapi.imrworldwide.com/")

	parser.WithMaxRequest(4)

	parser.ParseS3ObjectKey("EN660IH1TCIQT.2017-01-18-01.551a534d.gz",
		"us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod")

	log.Println("LOGPARSER STATS: ", parser.GetStats())

	parser = logparser.NewParser(
		"http://sandbox-cloudapi.imrworldwide.com/")

	parser.ParseS3ObjectKey("EN660IH1TCIQT.2017-01-18-01.551a534d.gz",
		"us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod")

	log.Println("LOGPARSER STATS: ", parser.GetStats())

	parser = logparser.NewParser(
		"http://sandbox-cloudapi.imrworldwide.com/")

	parser.ParseS3ObjectKey("EN660IH1TCIQT.2017-01-18-01.551a534d.gz",
		"us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod")

	log.Println("LOGPARSER STATS: ", parser.GetStats())

	parser = logparser.NewParser(
		"http://sandbox-cloudapi.imrworldwide.com/")

	parser.ParseS3ObjectKey("EN660IH1TCIQT.2017-01-18-01.551a534d.gz",
		"us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod")

	log.Println("LOGPARSER STATS: ", parser.GetStats())

}

func TestS3EventIntegration(*testing.T) {
	log.Println("TestIntegrationS3Event")
	eventData := []byte(`{
		"Records": [
		  {
			"eventVersion": "2.0",
			"eventTime": "1970-01-01T00:00:00.000Z",
			"requestParameters": {
			  "sourceIPAddress": "127.0.0.1"
			},
			"s3": {
			  "configurationId": "testConfigRule",
			  "object": {
				"eTag": "0123456789abcdef0123456789abcdef",
				"sequencer": "0A1B2C3D4E5F678901",
				"key": "EN660IH1TCIQT.2017-01-18-01.551a534d.gz",
				"size": 1024
			  },
			  "bucket": {
				"arn": "arn:aws:s3:::us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod",
				"name": "us-east-1-nlsn-data-dtvr-id3-aggregator-nonprod",
				"ownerIdentity": {
				  "principalId": "EXAMPLE"
				}
			  },
			  "s3SchemaVersion": "1.0"
			}
		  }
		]
	  }`)

	event := events.S3Event{}
	json.Unmarshal(eventData, &event)
	//log.Println(event)

	LambdaHandler(nil, event)

	log.Println("end")
}
