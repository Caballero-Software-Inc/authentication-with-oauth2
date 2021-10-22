package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TABLE_NAME = "users-caballero"

type Item struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Credits uint16 `json:"credits"`
}

func uploadUser(id string, email string) {
	aws_access_key_id := os.Getenv("AccessKeyID")
	aws_secret_access_key := os.Getenv("SecretAccessKey")
	region := os.Getenv("REGION")

	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, ""),
		Region:      &region,
	})

	svc := dynamodb.New(sess)

	inexisting := ExistingItem(id, email, svc)

	if inexisting {
		item := Item{
			Id:      id,
			Email:   email,
			Credits: 0,
		}
		AddDBItem(item, svc)
	} else {
		fmt.Println("The user was already registered.")
	}
}

func ExistingItem(id string, email string, svc *dynamodb.DynamoDB) bool {

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		fmt.Println("Error verifying item")
		fmt.Println(err.Error())
	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return (item.Id == "")
}

func AddDBItem(item Item, svc *dynamodb.DynamoDB) {
	av, _ := dynamodbattribute.MarshalMap(item)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TABLE_NAME),
	}
	_, err := svc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem: ")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("User registered")
}
