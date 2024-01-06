package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item"
	ErrorUserAlreadyExists       = "user.User already exists"
	ErrorUserDoesNotExist        = "user.User does not exist"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

func CreateUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (*User, error) {
	var user User

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	// validation level

	currentUser, _ := FetchUser(user.Email, awsTableName, dynamoDBClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	user.ID = uuid.New()
	attributeValue, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(awsTableName),
	}

	_, err = dynamoDBClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil
}

func FetchUser(
	email string,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(awsTableName),
	}

	result, err := dynamoDBClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchUsers(
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(awsTableName),
	}

	result, err := dynamoDBClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func UpdateUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) (*User, error) {
	var user User

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(user.Email, awsTableName, dynamoDBClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	attributeValue, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(awsTableName),
	}

	_, err = dynamoDBClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil
}

func DeleteUser(
	req events.APIGatewayProxyRequest,
	awsTableName string,
	dynamoDBClient dynamodbiface.DynamoDBAPI,
) error {
	email := req.QueryStringParameters["email"]

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(awsTableName),
	}

	_, err := dynamoDBClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
