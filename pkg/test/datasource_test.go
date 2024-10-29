package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/haohanyang/dynamodb-datasource/pkg/plugin"
)

func TestQueryData(t *testing.T) {
	ctx := context.Background()
	ds := plugin.CreateTestDatasource(ctx)

	err := createTable(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	err = writeItems(ctx, "test", []plugin.DataRow{
		{"value": &dynamodb.AttributeValue{
			N: aws.String("1"),
		}, "ts": &dynamodb.AttributeValue{
			N: aws.String("1730238174"),
		}},
		{"value": &dynamodb.AttributeValue{
			N: aws.String("2"),
		}, "ts": &dynamodb.AttributeValue{
			N: aws.String("1730238220"),
		}},
		{"value": &dynamodb.AttributeValue{
			N: aws.String("3"),
		}, "ts": &dynamodb.AttributeValue{
			N: aws.String("1730238227"),
		}},
	})

	if err != nil {
		t.Fatal(err)
	}

	qm := plugin.QueryModel{
		QueryText: "SELECT * FROM test WHERE id = 1 ORDER BY sid",
		Limit:     2,
		DatetimeFields: []plugin.DatetimeField{
			{
				Name:   "ts",
				Format: plugin.UnixTimestampSeconds,
			},
		},
	}

	rawJson, err := json.Marshal(qm)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: rawJson},
			},
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{},
				GrafanaConfig:              backend.NewGrafanaCfg(map[string]string{})},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}

	if resp.Responses["A"].Error != nil {
		t.Error(resp.Responses["A"].Error)
	} else {
		frame := resp.Responses["A"].Frames[0]
		size, err := frame.RowLen()
		if err != nil {
			t.Error(err)
		}
		assertEqual(t, size, 2)
		timeField, _ := frame.FieldByName("ts")
		valueField, _ := frame.FieldByName("value")

		ts1 := getFieldValue[time.Time](t, timeField, 0)
		assertEqual(t, ts1.Year(), 2024)
		assertEqual(t, ts1.Month(), time.October)
		assertEqual(t, ts1.Day(), 29)
		assertEqual(t, ts1.Hour(), 21)
		assertEqual(t, ts1.Minute(), 42)
		assertEqual(t, ts1.Second(), 54)

		ts2 := getFieldValue[time.Time](t, timeField, 1)
		assertEqual(t, ts2.Year(), 2024)
		assertEqual(t, ts2.Month(), time.October)
		assertEqual(t, ts2.Day(), 29)
		assertEqual(t, ts2.Hour(), 21)
		assertEqual(t, ts2.Minute(), 43)
		assertEqual(t, ts2.Second(), 40)

		assertEqual(t, valueField.At(0), plugin.Pointer[int64](1))
		assertEqual(t, valueField.At(1), plugin.Pointer[int64](2))
	}

}
