package server

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (s *server) validateTeamName(name string) (bool, error) {

	result, err := s.databaseClient.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.databaseClient.TeamDetailsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

func (s *server) validateHackathonName(name string) (bool, error) {

	result, err := s.databaseClient.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.databaseClient.HackathonDetailsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return false, nil
	}
	return true, nil

}
