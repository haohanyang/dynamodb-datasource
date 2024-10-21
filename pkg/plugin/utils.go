package plugin

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var endpoint = "http://localhost:4566"

func loadExtraPluginSettings(source backend.DataSourceInstanceSettings) (*ExtraPluginSettings, error) {
	settings := ExtraPluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	return &settings, nil
}

func parseNumber(n string) (*int64, *float64, error) {
	i, err := strconv.ParseInt(n, 10, 64)
	if err == nil {
		return aws.Int64(i), nil, nil
	} else {
		f, err := strconv.ParseFloat(n, 64)
		// float64
		if err == nil {
			return nil, aws.Float64(f), nil
		}
	}

	return nil, nil, fmt.Errorf("failed to parse %s", n)
}

func PrintDataFrame(dataFrame *data.Frame) {
	// Print headers
	fmt.Print("|")
	for i, field := range dataFrame.Fields {
		fmt.Print(field.Name)
		if i < len(dataFrame.Fields)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println("|")

	// Print data
	for i := 0; i < dataFrame.Rows(); i++ {
		fmt.Print("|")
		for j, field := range dataFrame.Fields {
			v, ok := field.ConcreteAt(i)

			if ok {
				if field.Type() == data.FieldTypeNullableJSON {
					rm := v.(json.RawMessage)
					rb, err := rm.MarshalJSON()
					if err != nil {
						panic(err)
					}
					fmt.Print(string(rb))
				} else if field.Type() == data.FieldTypeNullableString {
					s := v.(string)
					if len(s) > 10 {
						fmt.Print(s[:10] + "...")
					} else {
						fmt.Print(s)
					}
				} else {
					fmt.Print(v)
				}
			} else {
				fmt.Print("null")
			}

			if j < len(dataFrame.Fields)-1 {
				fmt.Print(",")
			}
		}
		fmt.Println("|")
	}
}

func OutputToDataFrame(dataFrameName string, output *dynamodb.ExecuteStatementOutput) (*data.Frame, error) {
	columns := make(map[string]*Column)
	for rowIndex, row := range output.Items {
		for name, value := range row {
			if c, ok := columns[name]; ok {
				err := c.AppendValue(value)
				if err != nil {
					return nil, err
				}
			} else {
				newColumn, err := NewColumn(rowIndex, name, value)
				if err != nil {
					return nil, err
				}
				if newColumn != nil {
					columns[name] = newColumn
				}
			}
		}

		// Make sure all columns have the same size
		for _, c := range columns {
			// Pad other columns with null value
			if c.Size() != rowIndex+1 {
				c.Field.Append(nil)
			}
		}
	}

	frame := data.NewFrame(dataFrameName)
	for _, c := range columns {
		frame.Fields = append(frame.Fields, c.Field)
	}

	return frame, nil
}

func NewTestClient() (*dynamodb.DynamoDB, error) {
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

func mapToJson(value *dynamodb.AttributeValue) (*json.RawMessage, error) {
	var m map[string]interface{}

	err := dynamodbattribute.UnmarshalMap(value.M, &m)
	if err != nil {
		return nil, err
	}

	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return pointer(json.RawMessage(jsonString)), nil
}

func listToJson(value *dynamodb.AttributeValue) (*json.RawMessage, error) {
	var l []interface{}

	err := dynamodbattribute.UnmarshalList(value.L, &l)
	if err != nil {
		return nil, err
	}

	jsonString, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return pointer(json.RawMessage(jsonString)), nil
}

func stringSetToJson(value *dynamodb.AttributeValue) (*json.RawMessage, error) {
	jsonString, err := json.Marshal(value.SS)
	if err != nil {
		return nil, err
	}
	return pointer(json.RawMessage(jsonString)), nil
}

func pointer[K any](val K) *K {
	return &val
}
