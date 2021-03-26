package server

import (
	"fmt"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/valyala/fasthttp"
)

func (s *server) validateHackathon() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present bool
		var err error
		name := string(ctx.QueryArgs().Peek("name"))
		if present, err = s.validateHackathonName(name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present {
			BasicResponse(400, "Hackathon name already used", ctx)
			return
		}

		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
	}
}

func (s *server) validateTeam() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present bool
		var err error
		name := string(ctx.QueryArgs().Peek("name"))
		if present, err = s.validateTeamName(name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present {
			BasicResponse(400, "Team name already used", ctx)
			return
		}

		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
	}
}

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
