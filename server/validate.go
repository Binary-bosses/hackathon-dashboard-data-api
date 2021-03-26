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
		var present map[string]*dynamodb.AttributeValue
		var err error
		name := string(ctx.QueryArgs().Peek("name"))
		if present, err = s.searchHackathonName(name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present != nil {
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
		var present map[string]*dynamodb.AttributeValue
		var err error
		name := string(ctx.QueryArgs().Peek("name"))
		if present, err = s.searchTeamName(name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present != nil {
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

func (s *server) validateHackathonAdmin() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present bool
		var err error
		name := string(ctx.QueryArgs().Peek("name"))
		pass := string(ctx.QueryArgs().Peek("pass"))
		if present, err = s.validateHackathonPass(name, pass); err != nil {
			BasicResponse(400, "Couldn't validate hackathon admin: "+err.Error(), ctx)
			return
		}
		if !present {
			BasicResponse(400, "Wrong pass, no edit access", ctx)
			return
		}

		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
	}
}

func (s *server) searchTeamName(name string) (map[string]*dynamodb.AttributeValue, error) {

	result, err := s.databaseClient.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.databaseClient.TeamDetailsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return nil, nil
	}
	return result.Item, nil
}

func (s *server) searchHackathonName(name string) (map[string]*dynamodb.AttributeValue, error) {

	result, err := s.databaseClient.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.databaseClient.HackathonDetailsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return nil, nil
	}
	return result.Item, nil

}

func (s *server) validateHackathonPass(name, pass string) (bool, error) {

	result, err := s.databaseClient.Service.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.databaseClient.HackathonPassTable),
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
	if password, ok := result.Item["pass"]; ok {
		if *password.S == pass {
			return true, nil
		}
	}
	return false, nil

}
