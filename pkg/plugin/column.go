package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type Column struct {
	Name  string
	Field *data.Field
}

func (c *Column) Type() data.FieldType {
	return c.Field.Type()
}

func NewColumn(rowIndex int, name string, value *dynamodb.AttributeValue) (*Column, error) {
	var field *data.Field

	if value.S != nil {
		// string
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, value.S)
	} else if value.N != nil {
		i, f, err := parseNumber(*value.N)
		if err != nil {
			return nil, err
		} else if i != nil {
			// int64
			field = data.NewField(name, nil, make([]*int64, rowIndex+1))
			field.Set(rowIndex, i)
		} else {
			// float64
			field = data.NewField(name, nil, make([]*float64, rowIndex+1))
			field.Set(rowIndex, f)
		}

	} else if value.B != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[B]"))
	} else if value.BOOL != nil {
		field = data.NewField(name, nil, make([]*bool, rowIndex+1))
		field.Set(rowIndex, value.BOOL)
	} else if value.NULL != nil {
		return nil, nil
	} else if value.M != nil {
		v, err := mapToJson(value)
		if err != nil {
			return nil, err
		}
		field = data.NewField(name, nil, make([]*json.RawMessage, rowIndex+1))
		field.Set(rowIndex, v)
	} else if value.L != nil {
		v, err := listToJson(value)
		if err != nil {
			return nil, err
		}
		field = data.NewField(name, nil, make([]*json.RawMessage, rowIndex+1))
		field.Set(rowIndex, v)
	} else if value.SS != nil {
		v, err := stringSetToJson(value)
		if err != nil {
			return nil, err
		}
		field = data.NewField(name, nil, make([]*json.RawMessage, rowIndex+1))
		field.Set(rowIndex, v)
	} else if value.NS != nil {
		v, err := numberSetToJson(value)
		if err != nil {
			return nil, err
		}
		field = data.NewField(name, nil, make([]*json.RawMessage, rowIndex+1))
		field.Set(rowIndex, v)
	} else if value.BS != nil {
		field = data.NewField(name, nil, make([]*string, rowIndex+1))
		field.Set(rowIndex, aws.String("[BS]"))
	}
	return &Column{Name: name, Field: field}, nil
}

func (c *Column) Size() int {
	return c.Field.Len()
}

func (c *Column) AppendValue(value *dynamodb.AttributeValue) error {
	if value.S != nil {
		if c.Type() != data.FieldTypeNullableString {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "S")
		}
		c.Field.Append(value.S)
	} else if value.N != nil {
		i, f, err := parseNumber(*value.N)
		if err != nil {
			return err
		} else if i != nil {
			// int64
			if c.Type() == data.FieldTypeNullableInt64 {
				c.Field.Append(i)
			} else if c.Type() == data.FieldTypeNullableFloat64 {
				c.Field.Append(aws.Float64(float64(*i)))
			} else {
				return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "N")
			}

		} else {
			// float64
			if c.Type() == data.FieldTypeNullableFloat64 {
				c.Field.Append(f)
			} else if c.Type() == data.FieldTypeNullableInt64 {

				// Convert all previous *int64 values to *float64
				float64Values := make([]*float64, c.Field.Len()+1)
				for i := 0; i < c.Field.Len(); i++ {
					cv, ok := c.Field.ConcreteAt(i)
					if ok {
						float64Values[i] = aws.Float64(float64(cv.(int64)))
					}
				}

				float64Values[c.Field.Len()] = f
				c.Field = data.NewField(c.Name, nil, float64Values)
			} else {
				return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "N")
			}
		}
	} else if value.B != nil {
		if c.Type() != data.FieldTypeNullableString {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "B")
		}
		c.Field.Append(aws.String("[B]"))
	} else if value.BOOL != nil {
		if c.Type() != data.FieldTypeNullableBool {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "BOOL")
		}
		c.Field.Append(value.BOOL)
	} else if value.NULL != nil {
		c.Field.Append(nil)
	} else if value.M != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "M")
		}
		v, err := mapToJson(value)
		if err != nil {
			return err
		}
		c.Field.Append(v)
	} else if value.L != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "L")
		}
		v, err := listToJson(value)
		if err != nil {
			return err
		}
		c.Field.Append(v)
	} else if value.SS != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "SS")
		}
		v, err := stringSetToJson(value)
		if err != nil {
			return err
		}
		c.Field.Append(v)
	} else if value.NS != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "NS")
		}
		v, err := numberSetToJson(value)
		if err != nil {
			return err
		}
		c.Field.Append(v)
	} else if value.BS != nil {
		if c.Type() != data.FieldTypeNullableString {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "BS")
		}
		c.Field.Append(aws.String("[BS]"))
	}

	return nil
}
