package plugin

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func parseNumber(n string) (*int64, *float64, error) {
	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return aws.Int64(i), nil, nil
	} else {
		f, err := strconv.ParseFloat(n, 64)
		// float64
		if err != nil {
			return nil, aws.Float64(f), nil
		}
	}

	return nil, nil, fmt.Errorf("failed to parse %s", n)
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
