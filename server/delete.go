package server

import (
	"errors"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/valyala/fasthttp"
)

func (s *server) deleteHackathon() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present map[string]*dynamodb.AttributeValue
		var err error
		name := string(ctx.QueryArgs().Peek("name"))
		if present, err = s.searchHackathonName(name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present == nil {
			BasicResponse(400, "Hackathon is not existing", ctx)
			return
		}
		var data HackathonData
		if err := dynamodbattribute.UnmarshalMap(present, &data); err != nil {
			BasicResponse(400, "Couldn't unmarshal team: "+err.Error(), ctx)
			return
		}

		for _, team := range data.Teams {

			if present, err = s.searchTeamName(team.Name); err != nil {
				BasicResponse(400, "Couldn't validate team name: "+err.Error(), ctx)
				return
			}
			if present == nil {
				BasicResponse(400, "team "+team.Name+" is not existing", ctx)
				return
			}

			if err := s.deleteTeam(team.Name); err != nil {
				BasicResponse(400, "couldn't delete team "+team.Name+" :"+err.Error(), ctx)
				return
			}
		}

		if err := s.deletePass(name); err != nil {
			BasicResponse(400, "couldn't delete hackathon pass"+err.Error(), ctx)
			return
		}

		if err := s.deleteHackDetails(name); err != nil {
			BasicResponse(400, "couldn't delete hackathon details"+err.Error(), ctx)
			return
		}

		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
	}
}

func (s *server) deleteHackDetails(name string) error {

	type itemDelete struct {
		Name string `json:"name"`
	}
	var team itemDelete
	team.Name = name

	av, err := dynamodbattribute.MarshalMap(team)
	if err != nil {
		return errors.New("Couldn't marshal team to delete")
	}

	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(s.databaseClient.HackathonDetailsTable),
	}

	_, err = s.databaseClient.Service.DeleteItem(input)
	if err != nil {
		return errors.New("Got error calling DeleteItem- " + err.Error())

	}

	return nil

}
func (s *server) deletePass(name string) error {

	type itemDelete struct {
		Name string `json:"name"`
	}
	var team itemDelete
	team.Name = name

	av, err := dynamodbattribute.MarshalMap(team)
	if err != nil {
		return errors.New("Couldn't marshal team to delete")
	}

	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(s.databaseClient.HackathonPassTable),
	}

	_, err = s.databaseClient.Service.DeleteItem(input)
	if err != nil {
		return errors.New("Got error calling DeleteItem- " + err.Error())

	}

	return nil

}
func (s *server) deleteTeam(name string) error {
	type itemDelete struct {
		Name string `json:"name"`
	}
	var team itemDelete
	team.Name = name

	av, err := dynamodbattribute.MarshalMap(team)
	if err != nil {
		return errors.New("Couldn't marshal team to delete")
	}

	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(s.databaseClient.TeamDetailsTable),
	}

	_, err = s.databaseClient.Service.DeleteItem(input)
	if err != nil {
		return errors.New("Got error calling DeleteItem- " + err.Error())

	}

	return nil
}
