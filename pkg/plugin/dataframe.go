package plugin

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func QueryResultToDataFrame(dataFrameName string, output *dynamodb.ExecuteStatementOutput, datetimeAttributes map[string]string) (*data.Frame, error) {
	attributes := make(map[string]*Attribute)
	for rowIndex, row := range output.Items {
		for name, value := range row {
			datetimeFormat := ""
			if df, ok := datetimeAttributes[name]; ok {
				datetimeFormat = df
			}

			if a, ok := attributes[name]; ok {
				err := a.Append(value)
				if err != nil {
					return nil, err
				}
			} else {
				newAttribute, err := NewAttribute(rowIndex, name, value, datetimeFormat)
				if err != nil {
					return nil, err
				}
				if newAttribute != nil {
					attributes[name] = newAttribute
				}
			}
		}

		// Make sure all attributes have the same size
		for _, c := range attributes {
			// Pad other attributes with null value
			if c.Size() != rowIndex+1 {
				c.Value.Append(nil)
			}
		}
	}

	frame := data.NewFrame(dataFrameName)
	for _, c := range attributes {
		frame.Fields = append(frame.Fields, c.Value)
	}

	return frame, nil
}
