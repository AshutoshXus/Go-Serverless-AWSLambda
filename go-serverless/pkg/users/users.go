package users

import (
	"encoding/json"
	"errors"

	"github.com/AshutoshXus/go-serverless/pkg/validators"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Creates Custom Error
var (
	ErorrFailedToFetchRecord     = "failed to fetch records"
	ErrorFailedToUnmarshalRecord = "failed to unmarshall the records"
	ErrorInvalidUserData         = "error invalid user data"
	ErrorInvalidEmail            = "error invalid email"
	ErrorMarshallData            = "error marshalling data"
	ErrorDeleteData              = "error deleting data"
	ErrorDynamoPutItem           = "could not dynamo put item"
	ErrorUserDoesnotExixts       = "error User doenot exists"
	ErrorUserDoesExixts          = "error User  exists"
)

func FetchUser(email string, tablename string, dynaclient dynamodbiface.DynamoDBAPI) (*User, error) {

	// Constructing the Query for DuynamoDB
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}

	// Geting the item based on Query
	result, err := dynaclient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErorrFailedToFetchRecord)
	}

	//Need to unmarshal the data from the Dynamodb so that frontend can understand the same.

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil

}

func FetchUsers(tablename string, dynaclient dynamodbiface.DynamoDBAPI) (*[]User, error) {

	// Constructing the Query for DuynamoDB
	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}

	// Geting the item based on Query
	result, err := dynaclient.Scan(input)

	if err != nil {
		return nil, errors.New(ErorrFailedToFetchRecord)
	}

	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return items, nil

}

func CreateUser(req events.APIGatewayProxyRequest, tablename string, dynaclient dynamodbiface.DynamoDBAPI) (*User, error) {

	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tablename, dynaclient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserDoesExixts)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorMarshallData)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = dynaclient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorDynamoPutItem)
	}

	return &u, nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dynaclient dynamodbiface.DynamoDBAPI) (*User, error) {

	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tablename, dynaclient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserDoesExixts)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorMarshallData)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = dynaclient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorDynamoPutItem)
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynaclient dynamodbiface.DynamoDBAPI) error {

	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}

	_, err := dynaclient.DeleteItem(input)

	if err != nil {
		return errors.New(ErrorDeleteData)
	}

	return nil

}
