package test

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/haohanyang/dynamodb-datasource/pkg/plugin"
)

var endpoint = "http://localhost:4566"
var testTableName = "test"

func testClient() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.AnonymousCredentials,
		Region:      aws.String("us-east-1"),
	})

	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}

func createTable(ctx context.Context, tableName string) error {
	client, err := testClient()
	if err != nil {
		return err
	}

	listTableOutput, err := client.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	hasTestTable := false
	for _, tn := range listTableOutput.TableNames {
		if *tn == tableName {
			hasTestTable = true
			break
		}
	}

	if hasTestTable {
		_, err = client.DeleteTableWithContext(ctx, &dynamodb.DeleteTableInput{
			TableName: &testTableName,
		})

		if err != nil {
			return err
		}
	}

	_, err = client.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("sid"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       aws.String("HASH"),
		}, {
			AttributeName: aws.String("sid"),
			KeyType:       aws.String("RANGE"),
		}},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	})

	return err
}

func writeItems(ctx context.Context, tableName string, items []plugin.DataRow) error {
	client, err := testClient()
	if err != nil {
		return err
	}

	for index, item := range items {
		c := maps.Clone(item)
		c["id"] = &dynamodb.AttributeValue{
			N: aws.String("1"),
		}
		c["sid"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.Itoa(index + 1)),
		}

		_, err := client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      c,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func outputFromItems(ctx context.Context, tableName string, items []plugin.DataRow, selectedAttributes string) (*dynamodb.ExecuteStatementOutput, error) {
	client, err := testClient()
	if err != nil {
		return nil, err
	}

	err = writeItems(ctx, tableName, items)
	if err != nil {
		return nil, err
	}

	return client.ExecuteStatementWithContext(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fmt.Sprintf("SELECT %s FROM %s WHERE id = 1", selectedAttributes, tableName)),
	})
}

func getFieldValue[T any](t *testing.T, field *data.Field, idx int) T {
	v, ok := field.ConcreteAt(idx)
	if !ok {
		t.Fatal("null plugin.Pointer")
	}
	return v.(T)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if cmp.Equal(a, b) {
		return
	}

	t.Errorf("Received %v (type %v), expected %v (type %v)", reflect.ValueOf(a), reflect.TypeOf(a), reflect.ValueOf(b), reflect.TypeOf(b))
}
