package server

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/valyala/fasthttp"
)

func (s *server) updateTeamDetails() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present map[string]*dynamodb.AttributeValue
		var err error
		oldTeam := Team{}
		createTeam := Team{}
		if err := json.Unmarshal(ctx.PostBody(), &createTeam); err != nil {
			BasicResponse(400, "Couldn't read request: "+err.Error(), ctx)
			return
		}
		if present, err = s.searchTeamName(createTeam.Name); err != nil {
			BasicResponse(400, "Couldn't validate team name '"+createTeam.Name+"' :"+err.Error(), ctx)
			return
		}
		if present == nil {
			BasicResponse(400, "Team name '"+createTeam.Name+"' not existing", ctx)
			return
		}

		if err := dynamodbattribute.UnmarshalMap(present, &oldTeam); err != nil {
			BasicResponse(400, "Couldn't unmarshal team: "+err.Error(), ctx)
			return
		}

		if err := s.updateTeam(oldTeam, createTeam); err != nil {
			BasicResponse(400, "Couldn't insert teams :"+err.Error(), ctx)
			return
		}
		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
		return

	}
}

func (s *server) registerHackathon() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present map[string]*dynamodb.AttributeValue
		var err error
		var editEvent HackathonData
		regEvent := EditHackathonData{}
		if err := json.Unmarshal(ctx.PostBody(), &regEvent); err != nil {
			BasicResponse(400, "Couldn't read request: "+err.Error(), ctx)
			return
		}
		if present, err = s.searchHackathonName(regEvent.Name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present == nil {
			BasicResponse(400, "Hackathon is not existing", ctx)
			return
		}

		if err := dynamodbattribute.UnmarshalMap(present, &editEvent); err != nil {
			BasicResponse(400, "Couldn't unmarshal team: "+err.Error(), ctx)
			return
		}

		for _, team := range regEvent.Teams {
			if present, err = s.searchTeamName(team.Name); err != nil {
				BasicResponse(400, "Couldn't validate team name: "+err.Error(), ctx)
				return
			}
			if present == nil {
				BasicResponse(400, "team not existing", ctx)
				return
			}
		}

		var hackathonData HackathonData

		hackathonData.Name = editEvent.Name
		hackathonData.StartTime = editEvent.StartTime
		hackathonData.EndTime = editEvent.EndTime
		hackathonData.Description = editEvent.Description
		hackathonData.Winner = editEvent.Winner
		for _, team := range editEvent.Teams {
			hackathonData.Teams = append(hackathonData.Teams, HackTeam{Name: team.Name, Idea: team.Idea})
		}

		for _, team := range regEvent.Teams {
			hackathonData.Teams = append(hackathonData.Teams, HackTeam{Name: team.Name, Idea: team.Idea})
		}

		if err := s.updateHackathonDetails(hackathonData, editEvent.Name); err != nil {
			BasicResponse(400, "couldn't update teams: "+err.Error(), ctx)
			return
		}
		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
	}
}

func (s *server) editHackathon() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		var present map[string]*dynamodb.AttributeValue
		var err error
		editEvent := EditHackathonData{}
		if err := json.Unmarshal(ctx.PostBody(), &editEvent); err != nil {
			BasicResponse(400, "Couldn't read request: "+err.Error(), ctx)
			return
		}

		if present, err = s.searchHackathonName(editEvent.Name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present == nil {
			BasicResponse(400, "Hackathon is not existing", ctx)
			return
		}

		for _, team := range editEvent.Teams {
			if present, err = s.searchTeamName(team.Name); err != nil {
				BasicResponse(400, "Couldn't validate team name: "+err.Error(), ctx)
				return
			}
			if present == nil {
				BasicResponse(400, "team is not existing", ctx)
				return
			}
		}

		hackathonData := HackathonData{}

		hackathonData.Name = editEvent.Name
		hackathonData.StartTime = editEvent.StartTime
		hackathonData.EndTime = editEvent.EndTime
		hackathonData.Description = editEvent.Description
		hackathonData.Winner = editEvent.Winner
		for _, team := range editEvent.Teams {
			hackathonData.Teams = append(hackathonData.Teams, HackTeam{Name: team.Name, Idea: team.Idea})
		}
		if err := s.updateHackathonDetails(hackathonData, editEvent.Name); err != nil {
			BasicResponse(400, "couldn't update teams: "+err.Error(), ctx)
			return
		}
		apiResp := APIResponse{
			Status: 200,
			Data:   Status{Status: "SUCCESS"},
		}

		util.SetJSONBody(ctx, apiResp)
		return
	}
}

func (s *server) updateHackathonDetails(data HackathonData, name string) error {

	var teams []*dynamodb.AttributeValue

	for _, team := range data.Teams {
		teams = append(teams, &dynamodb.AttributeValue{M: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(team.Name),
			},
			"idea": {
				S: aws.String(team.Idea),
			},
		},
		})
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":teams": {
				L: teams,
			},
			":description": {
				S: aws.String(data.Description),
			},
			":startTime": {
				N: aws.String(fmt.Sprintf("%f", data.StartTime.(float64))),
			},
			":endTime": {
				N: aws.String(fmt.Sprintf("%f", data.EndTime.(float64))),
			},
		},
		TableName: aws.String(s.databaseClient.HackathonDetailsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set teams = :teams, description = :description, startTime = :startTime, endTime = :endTime "),
	}

	if data.Winner != "" {
		input.ExpressionAttributeValues[":winner"] = &dynamodb.AttributeValue{S: aws.String(data.Winner)}
		input.UpdateExpression = aws.String("set teams = :teams, description = :description, startTime = :startTime, endTime = :endTime, winner= :winner")
	}

	_, err := s.databaseClient.Service.UpdateItem(input)
	if err != nil {
		return errors.New("Got error calling UpdateItem: " + err.Error())
	}

	return nil
}
func (s *server) updateTeam(oldteam, newTeam Team) error {

	var members []*dynamodb.AttributeValue

	for _, mem := range newTeam.Members {
		members = append(members, &dynamodb.AttributeValue{S: aws.String(mem)})
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":members": {
				L: members,
			},
		},
		TableName: aws.String(s.databaseClient.TeamDetailsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(newTeam.Name),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set members = :members"),
	}

	_, err := s.databaseClient.Service.UpdateItem(input)
	if err != nil {
		return errors.New("Got error calling UpdateItem: " + err.Error())
	}

	return nil
}
