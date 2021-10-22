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

type Item struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Credits uint16 `json:"credits"`
}

func uploadUser(id string, email string) {
	fmt.Println(email)
	fmt.Println(id)
	item := Item{
		Id:      id,
		Email:   email,
		Credits: 0,
	}
	//send it to AWS DynamoDB

	aws_access_key_id := os.Getenv("AccessKeyID")
	aws_secret_access_key := os.Getenv("SecretAccessKey")
	region := os.Getenv("REGION")

	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, ""),
		Region:      &region,
	})

	svc := dynamodb.New(sess)
	AddDBItem(item, svc)
}

func AddDBItem(item Item, svc *dynamodb.DynamoDB) {
	TABLE_NAME := "users-caballero"
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

	fmt.Println("Item added to the table.")
}
