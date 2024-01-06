package handlers

import (
	"aws_integration_with_golang/pkg/user"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func CreateUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, awsTableName, dynamoDBClient)
	if err != nil {
		awsFormatError := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusMethodNotAllowed, awsFormatError)
	}

	return apiResponse(http.StatusCreated, result)
}

func GetUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, awsTableName, dynamoDBClient)
		if err != nil {
			awsFormatError := ErrorBody{aws.String(err.Error())}
			return apiResponse(http.StatusBadRequest, awsFormatError)
		}

		return apiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(awsTableName, dynamoDBClient)
	if err != nil {
		awsFormatError := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, awsFormatError)
	}

	return apiResponse(http.StatusOK, result)
}

func UpdateUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, awsTableName, dynamoDBClient)
	if err != nil {
		awsFormatError := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, awsFormatError)
	}

	return apiResponse(http.StatusOK, result)
}

func DeleteUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, awsTableName, dynamoDBClient)
	if err != nil {
		awsFormatError := ErrorBody{aws.String(err.Error())}
		return apiResponse(http.StatusBadRequest, awsFormatError)
	}

	return apiResponse(http.StatusOK, nil)
}

func UnhadleMethod() (events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
