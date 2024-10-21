package plugin

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewTestClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String("http://localhost:4566"),
		Credentials: credentials.AnonymousCredentials,
		Region:      aws.String("us-east-1"),
	})

	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}
