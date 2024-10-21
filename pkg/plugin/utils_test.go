package plugin

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var testTableName = "test"

func newTestClient() (*dynamodb.DynamoDB, error) {
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

func writeItems(ctx context.Context, client *dynamodb.DynamoDB, tableName string, items []DataRow) error {
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

func outputFromItems(ctx context.Context, client *dynamodb.DynamoDB, tableName string, items []DataRow, selectedAttributes string) (*dynamodb.ExecuteStatementOutput, error) {
	err := writeItems(ctx, client, tableName, items)
	if err != nil {
		return nil, err
	}

	return client.ExecuteStatementWithContext(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fmt.Sprintf("SELECT %s from %s where id = 1", selectedAttributes, tableName)),
	})
}

func getFieldValue[T any](t *testing.T, field *data.Field, idx int) T {
	v, ok := field.ConcreteAt(idx)
	if !ok {
		t.Fatal("null pointer")
	}
	return v.(T)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}

	t.Errorf("Received %v (type %v), expected %v (type %v)", reflect.ValueOf(a), reflect.TypeOf(a), reflect.ValueOf(b), reflect.TypeOf(b))
}

func TestOutputToDataFrame(t *testing.T) {
	ctx := context.Background()
	client, err := newTestClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.DeleteTableWithContext(ctx, &dynamodb.DeleteTableInput{
		TableName: &testTableName,
	})

	if err != nil && !strings.Contains(err.Error(), "ResourceNotFoundException") {
		t.Fatal(err)
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

	if err != nil {
		t.Fatal(err)
	}

	t.Run("N int", func(t *testing.T) {
		rows := []DataRow{
			{"myNI": &dynamodb.AttributeValue{
				N: aws.String("1"),
			}},
			{"my": &dynamodb.AttributeValue{
				N: aws.String("2"),
			}},
			{"myNI": &dynamodb.AttributeValue{
				N: aws.String("2"),
			}},
		}

		output, err := outputFromItems(ctx, client, testTableName, rows, "myNI")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := OutputToDataFrame("test", output)
		if err != nil {
			t.Fatal(err)
		}

		field := frame.Fields[0]
		var null *int64
		assertEqual(t, field.Name, "myNI")
		assertEqual(t, getFieldValue[int64](t, field, 0), *aws.Int64(1))
		assertEqual(t, field.At(1), null)
		assertEqual(t, getFieldValue[int64](t, field, 2), *aws.Int64(2))
	})

	t.Run("N float", func(t *testing.T) {
		rows := []DataRow{
			{"myNF": &dynamodb.AttributeValue{
				N: aws.String("1.2"),
			}},
			{"my": &dynamodb.AttributeValue{
				N: aws.String("2"),
			}},
			{"myNF": &dynamodb.AttributeValue{
				N: aws.String("2.1"),
			}},
		}

		output, err := outputFromItems(ctx, client, testTableName, rows, "myNF")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := OutputToDataFrame("test", output)
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		var null *float64
		assertEqual(t, field.Name, "myNF")
		assertEqual(t, fmt.Sprintf("%.1f", getFieldValue[float64](t, field, 0)), "1.2")
		assertEqual(t, field.At(1), null)
		assertEqual(t, fmt.Sprintf("%.1f", getFieldValue[float64](t, field, 2)), "2.1")
	})

	t.Run("N int & float", func(t *testing.T) {
		rows := []DataRow{
			{"myNIF": &dynamodb.AttributeValue{
				N: aws.String("1"),
			}},
			{"my": &dynamodb.AttributeValue{
				N: aws.String("2"),
			}},
			{"myNIF": &dynamodb.AttributeValue{
				N: aws.String("2.1"),
			}},
		}

		output, err := outputFromItems(ctx, client, testTableName, rows, "myNIF")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := OutputToDataFrame("test", output)
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		var null *float64
		assertEqual(t, field.Name, "myNIF")
		assertEqual(t, fmt.Sprintf("%.1f", getFieldValue[float64](t, field, 0)), "1.0")
		assertEqual(t, field.At(1), null)
		assertEqual(t, fmt.Sprintf("%.1f", getFieldValue[float64](t, field, 2)), "2.1")
	})

	t.Run("N float & int", func(t *testing.T) {
		rows := []DataRow{
			{"myNFI": &dynamodb.AttributeValue{
				N: aws.String("1.1"),
			}},
			{"my": &dynamodb.AttributeValue{
				N: aws.String("2"),
			}},
			{"myNFI": &dynamodb.AttributeValue{
				N: aws.String("2"),
			}},
		}

		output, err := outputFromItems(ctx, client, testTableName, rows, "myNFI")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := OutputToDataFrame("test", output)
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		var null *float64
		assertEqual(t, field.Name, "myNFI")
		assertEqual(t, fmt.Sprintf("%.1f", getFieldValue[float64](t, field, 0)), "1.1")
		assertEqual(t, field.At(1), null)
		assertEqual(t, fmt.Sprintf("%.1f", getFieldValue[float64](t, field, 2)), "2.0")
	})

	t.Run("BOOL", func(t *testing.T) {
		rows := []DataRow{
			{"myBOOL": &dynamodb.AttributeValue{
				BOOL: aws.Bool(true),
			}},
			{"myBOOL": &dynamodb.AttributeValue{
				BOOL: aws.Bool(false),
			}},
		}

		output, err := outputFromItems(ctx, client, testTableName, rows, "myBOOL")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := OutputToDataFrame("test", output)
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myBOOL")
		assertEqual(t, getFieldValue[bool](t, field, 0), true)
		assertEqual(t, getFieldValue[bool](t, field, 1), false)
	})

	t.Run("M", func(t *testing.T) {
		rows := []DataRow{
			{"myM": &dynamodb.AttributeValue{
				M: map[string]*dynamodb.AttributeValue{
					"key1": {
						S: aws.String("string1"),
					},
					"key2": {
						N: aws.String("1"),
					},
				},
			}},
			{"myM": &dynamodb.AttributeValue{
				M: map[string]*dynamodb.AttributeValue{
					"key3": {
						S: aws.String("string2"),
					},
					"key4": {
						N: aws.String("2.1"),
					},
				},
			}},
		}

		output, err := outputFromItems(ctx, client, testTableName, rows, "myM")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := OutputToDataFrame("test", output)
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myM")
		assertEqual(t, field.Type(), data.FieldTypeNullableJSON)
	})

}
