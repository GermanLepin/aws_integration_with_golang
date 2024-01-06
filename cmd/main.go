package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	awsTableName = "LambdaInGoUser"
)

var (
	dynamoDBClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")

	awsConfig := aws.Config{
		Region: aws.String(region),
	}

	awsSession, err := session.NewSession(&awsConfig)
	if err != nil {
		return
	}

	dynamoDBClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodGet:
		return handlers.GetUser(req, awsTableName, dynamoDBClient)
	case http.MethodPost:
		return handlers.CreateUser(req, awsTableName, dynamoDBClient)
	case http.MethodPut:
		return handlers.UpdateUser(req, awsTableName, dynamoDBClient)
	case http.MethodDelete:
		return handlers.DeleteUser(req, awsTableName, dynamoDBClient)
	default:
		return handlers.UnhadleMethod()
	}
}
