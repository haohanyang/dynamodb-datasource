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

	t.Run("unix datetime", func(t *testing.T) {
		err = writeItems(ctx, "test", []plugin.DataRow{
			{"ts": &dynamodb.AttributeValue{
				N: aws.String("1730238174"),
			}},
			{"ts": &dynamodb.AttributeValue{
				N: aws.String("1730324262"),
			}},
		})

		if err != nil {
			t.Fatal(err)
		}

		qm := plugin.QueryModel{
			QueryText: "SELECT * FROM test",
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

			ts1 := getFieldValue[time.Time](t, timeField, 0)
			ts2 := getFieldValue[time.Time](t, timeField, 1)
			assertEqual(t, ts1.Unix(), int64(1730238174))
			assertEqual(t, ts2.Unix(), int64(1730324262))
		}
	})

	t.Run("iso datetime", func(t *testing.T) {
		err = writeItems(ctx, "test", []plugin.DataRow{
			{"ts": &dynamodb.AttributeValue{
				S: aws.String("2024-10-30T22:46:09+01:00"),
			}},
			{"ts": &dynamodb.AttributeValue{
				S: aws.String("2024-10-30T22:46:34+01:00"),
			}},
		})

		if err != nil {
			t.Fatal(err)
		}

		qm := plugin.QueryModel{
			QueryText: "SELECT * FROM test",
			Limit:     2,
			DatetimeFields: []plugin.DatetimeField{
				{
					Name:   "ts",
					Format: "2006-01-02T15:04:05-07:00",
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

			ts1 := getFieldValue[time.Time](t, timeField, 0)
			ts2 := getFieldValue[time.Time](t, timeField, 1)
			assertEqual(t, ts1.Unix(), int64(1730324769))
			assertEqual(t, ts2.Unix(), int64(1730324794))
		}
	})
}
