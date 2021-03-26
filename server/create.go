package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/valyala/fasthttp"
)

func (s *server) createHackathon() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var present bool
		var err error
		createEvent := CreateHackathonData{}
		if err := json.Unmarshal(ctx.PostBody(), &createEvent); err != nil {
			BasicResponse(400, "Couldn't read request: "+err.Error(), ctx)
			return
		}

		if present, err = s.validateHackathonName(createEvent.Name); err != nil {
			BasicResponse(400, "Couldn't validate hackathon name: "+err.Error(), ctx)
			return
		}
		if present {
			BasicResponse(400, "Hackathon name already used", ctx)
			return
		}

		for _, teams := range createEvent.Teams {
			if present, err = s.validateTeamName(teams.Name); err != nil {
				BasicResponse(400, "Couldn't validate team name '"+teams.Name+"' :"+err.Error(), ctx)
				return
			}
			if present {
				BasicResponse(400, "Team name '"+teams.Name+"' already used", ctx)
				return
			}

		}

		if err := s.insertTeams(createEvent.Teams); err != nil {
			BasicResponse(400, "Couldn't insert teams :"+err.Error(), ctx)
			return
		}

		if err := s.insertEditPass(createEvent.Name, createEvent.HackathonPass); err != nil {
			BasicResponse(400, "Couldn't insert hackathon pass :"+err.Error(), ctx)
			return
		}

		if err := s.insertHackathonDetails(createEvent); err != nil {
			BasicResponse(400, "Couldn't insert hackathon pass :"+err.Error(), ctx)
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

func (s *server) insertTeams(teams []Team) error {
	for _, team := range teams {
		av, err := dynamodbattribute.MarshalMap(team)
		if err != nil {
			return errors.New("Got error marshalling new team item: " + err.Error())
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(s.databaseClient.TeamDetailsTable),
		}

		_, err = s.databaseClient.Service.PutItem(input)
		if err != nil {
			return errors.New("Got error calling PutItem: " + err.Error())
		}

		log.Println("Inserted team :" + team.Name)
	}
	return nil
}

func (s *server) insertEditPass(hackName, passName string) error {
	item := HackathonEditPass{}

	item.Name = hackName
	item.Pass = passName
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return errors.New("Got error marshalling new hackthon edit pass item: " + err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.databaseClient.HackathonPassTable),
	}

	_, err = s.databaseClient.Service.PutItem(input)
	if err != nil {
		return errors.New("Got error calling PutItem: " + err.Error())
	}

	log.Println("Inserted pass for hackathon:" + item.Name)
	return nil
}

func (s *server) insertHackathonDetails(data CreateHackathonData) error {
	hackathonData := HackathonData{}

	hackathonData.Name = data.Name
	hackathonData.StartTime = data.StartTime
	hackathonData.EndTime = data.EndTime
	hackathonData.Description = data.Description
	for _, team := range data.Teams {
		hackathonData.Teams = append(hackathonData.Teams, team.Name)
	}

	av, err := dynamodbattribute.MarshalMap(hackathonData)
	if err != nil {
		return errors.New("Got error marshalling new hackathon item: " + err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.databaseClient.HackathonDetailsTable),
	}

	_, err = s.databaseClient.Service.PutItem(input)
	if err != nil {
		return errors.New("Got error calling PutItem: " + err.Error())
	}

	log.Println("Inserted details of hackathon:" + hackathonData.Name)
	return nil

}
