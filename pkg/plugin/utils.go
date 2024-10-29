package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func loadExtraPluginSettings(source backend.DataSourceInstanceSettings) (*ExtraPluginSettings, error) {
	settings := ExtraPluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	return &settings, nil
}

func CreateTestDatasource(ctx context.Context) *Datasource {
	dsSetting := awsds.AWSDatasourceSettings{
		Profile:   "test",
		Region:    "us-east-1",
		AuthType:  awsds.AuthTypeKeys,
		Endpoint:  "http://localhost:4566",
		AccessKey: "test",
		SecretKey: "test",
	}

	authSettings := awsds.ReadAuthSettings(ctx)
	sessionCache := awsds.NewSessionCache()

	return &Datasource{
		Settings:     dsSetting,
		authSettings: *authSettings,
		sessionCache: sessionCache,
	}
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
	return Pointer(json.RawMessage(jsonString)), nil
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
	return Pointer(json.RawMessage(jsonString)), nil
}

func stringSetToJson(value *dynamodb.AttributeValue) (*json.RawMessage, error) {
	jsonString, err := json.Marshal(value.SS)
	if err != nil {
		return nil, err
	}
	return Pointer(json.RawMessage(jsonString)), nil
}

func numberSetToJson(value *dynamodb.AttributeValue) (*json.RawMessage, error) {
	l := make([]interface{}, len(value.NS))

	for idx, n := range value.NS {
		i, f, err := parseNumber(*n)
		if err != nil {
			return nil, err
		}
		if i != nil {
			l[idx] = *i
		} else {
			l[idx] = *f
		}
	}

	jsonString, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return Pointer(json.RawMessage(jsonString)), nil
}

func Pointer[K any](val K) *K {
	return &val
}
