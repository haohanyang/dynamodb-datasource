package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/haohanyang/dynamodb-datasource/pkg/plugin"
)

func TestQueryResultToDataFrame(t *testing.T) {
	ctx := context.Background()
	err := createTable(ctx, "test")

	if err != nil {
		t.Fatal(err)
	}
	t.Run("N int", func(t *testing.T) {
		rows := []plugin.DataRow{
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

		output, err := outputFromItems(ctx, testTableName, rows, "myNI")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
		if err != nil {
			t.Fatal(err)
		}

		field := frame.Fields[0]
		var null *int64
		assertEqual(t, field.Name, "myNI")
		assertEqual(t, field.At(0), plugin.Pointer[int64](1))
		assertEqual(t, field.At(1), null)
		assertEqual(t, field.At(2), plugin.Pointer[int64](2))
	})

	t.Run("N float", func(t *testing.T) {
		rows := []plugin.DataRow{
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

		output, err := outputFromItems(ctx, testTableName, rows, "myNF")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
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
		rows := []plugin.DataRow{
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

		output, err := outputFromItems(ctx, testTableName, rows, "myNIF")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
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
		rows := []plugin.DataRow{
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

		output, err := outputFromItems(ctx, testTableName, rows, "myNFI")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
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
		rows := []plugin.DataRow{
			{"myBOOL": &dynamodb.AttributeValue{
				BOOL: aws.Bool(true),
			}},
			{"myBOOL": &dynamodb.AttributeValue{
				BOOL: aws.Bool(false),
			}},
		}

		output, err := outputFromItems(ctx, testTableName, rows, "myBOOL")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myBOOL")
		assertEqual(t, field.At(0), plugin.Pointer(true))
		assertEqual(t, field.At(1), plugin.Pointer(false))
	})

	t.Run("M", func(t *testing.T) {
		rows := []plugin.DataRow{
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

		output, err := outputFromItems(ctx, testTableName, rows, "myM")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myM")
		assertEqual(t, field.Type(), data.FieldTypeNullableJSON)
	})

	t.Run("L", func(t *testing.T) {
		rows := []plugin.DataRow{
			{"myL": &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{
					{
						BOOL: aws.Bool(true),
					},
					{
						N: aws.String("1"),
					},
				},
			}},
			{"myL": &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{
					{
						S: aws.String("string2"),
					},
					{
						N: aws.String("2.1"),
					},
				},
			}},
		}

		output, err := outputFromItems(ctx, testTableName, rows, "myL")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myL")
		assertEqual(t, field.Type(), data.FieldTypeNullableJSON)
	})

	t.Run("SS", func(t *testing.T) {
		rows := []plugin.DataRow{
			{"mySS": &dynamodb.AttributeValue{
				SS: []*string{aws.String("s1"), aws.String("s2")},
			}},
			{"mySS": &dynamodb.AttributeValue{
				SS: []*string{aws.String("s3"), aws.String("s4")},
			}},
		}

		output, err := outputFromItems(ctx, testTableName, rows, "mySS")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "mySS")
		assertEqual(t, field.Type(), data.FieldTypeNullableJSON)
	})

	t.Run("NS", func(t *testing.T) {
		rows := []plugin.DataRow{
			{"myNS": &dynamodb.AttributeValue{
				NS: []*string{aws.String("1.1"), aws.String("2")},
			}},
			{"myNS": &dynamodb.AttributeValue{
				NS: []*string{aws.String("-2"), aws.String("-3.1")},
			}},
		}

		output, err := outputFromItems(ctx, testTableName, rows, "myNS")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, make(map[string]string))
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myNS")
		assertEqual(t, field.Type(), data.FieldTypeNullableJSON)
	})

	t.Run("Datetime Unix seconds", func(t *testing.T) {
		rows := []plugin.DataRow{
			{"myDate": &dynamodb.AttributeValue{
				N: aws.String("1730070176"),
			}},
			{},
			{"myDate": &dynamodb.AttributeValue{
				N: aws.String("1730070193"),
			}}}

		output, err := outputFromItems(ctx, testTableName, rows, "myDate")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, map[string]string{
			"myDate": plugin.UnixTimestampSeconds})
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myDate")
		assertEqual(t, field.Type(), data.FieldTypeNullableTime)
		v, _ := field.ConcreteAt(0)
		dt := v.(time.Time)
		assertEqual(t, dt.Year(), 2024)
		assertEqual(t, dt.Month().String(), "October")
		assertEqual(t, dt.Day(), 27)
	})

	t.Run("Datetime Unix miliseconds", func(t *testing.T) {
		rows := []plugin.DataRow{
			{"myDate": &dynamodb.AttributeValue{
				N: aws.String("1730070554000"),
			}},
			{},
			{"myDate": &dynamodb.AttributeValue{
				N: aws.String("1730070568000"),
			}}}

		output, err := outputFromItems(ctx, testTableName, rows, "myDate")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, map[string]string{
			"myDate": plugin.UnixTimestampMiniseconds})
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myDate")
		assertEqual(t, field.Type(), data.FieldTypeNullableTime)
		v, _ := field.ConcreteAt(0)
		dt := v.(time.Time)
		assertEqual(t, dt.Year(), 2024)
		assertEqual(t, dt.Month().String(), "October")
		assertEqual(t, dt.Day(), 27)
	})

	t.Run("Datetime Custom format ISO8601", func(t *testing.T) {

		rows := []plugin.DataRow{
			{"myDate": &dynamodb.AttributeValue{
				S: aws.String("2024-10-27T23:10:42.951Z"),
			}},
			{},
			{"myDate": &dynamodb.AttributeValue{
				S: aws.String("2024-10-27T23:10:49.552Z"),
			}}}

		output, err := outputFromItems(ctx, testTableName, rows, "myDate")
		if err != nil {
			t.Fatal(err)
		}

		frame, err := plugin.QueryResultToDataFrame("test", output, map[string]string{
			"myDate": "2006-01-02T15:04:05.999Z"})
		if err != nil {
			t.Fatal(err)
		}
		field := frame.Fields[0]
		assertEqual(t, field.Name, "myDate")
		assertEqual(t, field.Type(), data.FieldTypeNullableTime)
		v, _ := field.ConcreteAt(0)
		dt := v.(time.Time)
		assertEqual(t, dt.Year(), 2024)
		assertEqual(t, dt.Month().String(), "October")
		assertEqual(t, dt.Day(), 27)
		assertEqual(t, dt.Hour(), 23)
		assertEqual(t, dt.Minute(), 10)
		assertEqual(t, dt.Second(), 42)
	})

}
