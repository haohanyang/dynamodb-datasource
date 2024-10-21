package plugin

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type Column struct {
	Name  string
	Field *data.Field
}

func NewColumn(rowIndex int, name string, value *dynamodb.AttributeValue) *Column {
	var field *data.Field

	if value.S != nil {
		// string
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, value.S)
	} else if value.N != nil {
		// int64
		i, err := strconv.ParseInt(*value.N, 10, 64)
		if err != nil {
			field = data.NewField(name, nil, make([]*int64, rowIndex+1))
			field.Set(rowIndex, i)
		} else {
			f, err := strconv.ParseFloat(*value.N, 64)
			// float64
			if err != nil {
				field = data.NewField(name, nil, make([]*float64, rowIndex+1))
				field.Set(rowIndex, f)
			} else {
				backend.Logger.Error("Failed to parse N", *value.N)
			}
		}
	} else if value.B != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[BINARY]"))
	} else if value.BOOL != nil {
		field = data.NewField(name, nil, make([]*bool, rowIndex+1))
		field.Set(rowIndex, value.BOOL)
	} else if value.M != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[MAP]"))
	} else if value.L != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[LIST]"))
	} else if value.SS != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[SS]"))
	} else if value.NS != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[NS]"))
	} else if value.BS != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[BS]"))
	}
	return nil
}
