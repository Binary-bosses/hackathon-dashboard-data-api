package database

import "github.com/aws/aws-sdk-go/service/dynamodb"

type Database struct {
	Region                string
	Service               *dynamodb.DynamoDB
	HackathonDetailsTable string
	TeamDetailsTable      string
	HackathonPassTable    string
}
