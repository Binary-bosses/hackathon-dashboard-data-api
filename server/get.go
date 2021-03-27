package server

import (
	"time"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/valyala/fasthttp"
)

func (s *server) getHackathons() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		status := string(ctx.QueryArgs().Peek("status"))
		var filter expression.ConditionBuilder
		var expr expression.Expression
		var err error
		if status == "completed" {
			filter = expression.Name("endTime").LessThan(expression.Value(float64(time.Now().Unix()))).And(expression.Name("startTime").LessThan(expression.Value(float64(time.Now().Unix()))))
		} else if status == "ongoing" {
			filter = expression.Name("endTime").GreaterThan(expression.Value(float64(time.Now().Unix()))).And(expression.Name("startTime").LessThan(expression.Value(float64(time.Now().Unix()))))
		} else if status == "upcoming" {
			filter = expression.Name("endTime").GreaterThan(expression.Value(float64(time.Now().Unix()))).And(expression.Name("startTime").GreaterThan(expression.Value(float64(time.Now().Unix()))))
		}

		proj := expression.NamesList(expression.Name("name"), expression.Name("description"), expression.Name("startTime"), expression.Name("endTime"))
		if status != "" {
			expr, err = expression.NewBuilder().WithFilter(filter).WithProjection(proj).Build()
			if err != nil {
				BasicResponse(400, "Got error building expression:"+err.Error(), ctx)
				return
			}
		} else {
			expr, err = expression.NewBuilder().WithProjection(proj).Build()
			if err != nil {
				BasicResponse(400, "Got error building expression:"+err.Error(), ctx)
				return
			}
		}

		// Build the query input parameters
		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String(s.databaseClient.HackathonDetailsTable),
		}

		// Make the DynamoDB Query API call
		result, err := s.databaseClient.Service.Scan(params)
		if err != nil {
			BasicResponse(400, "Query failed :"+err.Error(), ctx)
			return
		}

		var items []HackathonHighLevelData
		for _, i := range result.Items {
			item := HackathonHighLevelData{}

			err = dynamodbattribute.UnmarshalMap(i, &item)

			if err != nil {
				BasicResponse(400, "Couldn't unmarshal hackahon data: "+err.Error(), ctx)
				return
			}

			items = append(items, item)
		}

		apiResp := APIResponse{
			Status: 200,
			Data:   items,
		}

		util.SetJSONBody(ctx, apiResp)

	}

}

func (s *server) getHackathon() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		name := string(ctx.QueryArgs().Peek("name"))

		filter := expression.Name("name").Equal(expression.Value(name))

		proj := expression.NamesList(expression.Name("name"), expression.Name("description"), expression.Name("startTime"), expression.Name("endTime"), expression.Name("teams"), expression.Name("winner"))

		expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(proj).Build()
		if err != nil {
			BasicResponse(400, "Got error building expression:"+err.Error(), ctx)
			return
		}

		// Build the query input parameters
		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String(s.databaseClient.HackathonDetailsTable),
		}

		// Make the DynamoDB Query API call
		result, err := s.databaseClient.Service.Scan(params)
		if err != nil {
			BasicResponse(400, "Query failed :"+err.Error(), ctx)
			return
		}

		for _, i := range result.Items {
			item := HackathonData{}

			err = dynamodbattribute.UnmarshalMap(i, &item)

			if err != nil {
				BasicResponse(400, "Couldn't unmarshal hackahon data: "+err.Error(), ctx)
				return
			}

			apiResp := APIResponse{
				Status: 200,
				Data:   item,
			}

			util.SetJSONBody(ctx, apiResp)

			return
		}

		apiResp := APIResponse{
			Status: 200,
		}

		util.SetJSONBody(ctx, apiResp)

		return
	}

}

func (s *server) getTeam() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		name := string(ctx.QueryArgs().Peek("name"))

		filter := expression.Name("name").Equal(expression.Value(name))

		proj := expression.NamesList(expression.Name("name"), expression.Name("members"))

		expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(proj).Build()
		if err != nil {
			BasicResponse(400, "Got error building expression:"+err.Error(), ctx)
			return
		}

		// Build the query input parameters
		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String(s.databaseClient.TeamDetailsTable),
		}

		// Make the DynamoDB Query API call
		result, err := s.databaseClient.Service.Scan(params)
		if err != nil {
			BasicResponse(400, "Query failed :"+err.Error(), ctx)
			return
		}

		for _, i := range result.Items {
			item := Team{}

			err = dynamodbattribute.UnmarshalMap(i, &item)

			if err != nil {
				BasicResponse(400, "Couldn't unmarshal hackahon data: "+err.Error(), ctx)
				return
			}

			apiResp := APIResponse{
				Status: 200,
				Data:   item,
			}

			util.SetJSONBody(ctx, apiResp)

			return
		}

		apiResp := APIResponse{
			Status: 200,
		}

		util.SetJSONBody(ctx, apiResp)

		return

	}
}
