package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type Attribute struct {
	Name     string
	Value    *data.Field
	TsFormat string
}

func (c *Attribute) Type() data.FieldType {
	return c.Value.Type()
}

func NewAttribute(rowIndex int, name string, value *dynamodb.AttributeValue, datetimeFormat string) (*Attribute, error) {
	var field *data.Field

	if value.S != nil {
		// string
		if datetimeFormat != "" && datetimeFormat != UnixTimestampMiniseconds && datetimeFormat != UnixTimestampSeconds {
			t, err := time.Parse(datetimeFormat, *value.S)
			if err != nil {
				return nil, err
			}

			field = data.NewField(name, nil, make([]*time.Time, rowIndex+1))
			field.Set(rowIndex, &t)

		} else {
			field = data.NewField(name, nil, make([]*string, rowIndex+1))
			field.Set(rowIndex, value.S)
		}

	} else if value.N != nil {
		i, f, err := parseNumber(*value.N)
		if err != nil {
			return nil, err
		} else if i != nil {
			// int64
			if datetimeFormat == UnixTimestampSeconds {
				t := time.Unix(*i, 0)
				field = data.NewField(name, nil, make([]*time.Time, rowIndex+1))
				field.Set(rowIndex, &t)
			} else if datetimeFormat == UnixTimestampMiniseconds {
				seconds := *i / 1000
				nanoseconds := (*i % 1000) * 1000000
				t := time.Unix(seconds, nanoseconds)
				field = data.NewField(name, nil, make([]*time.Time, rowIndex+1))
				field.Set(rowIndex, &t)
			} else if datetimeFormat != "" {
				return nil, errors.New("invalid datetime format")
			} else {
				field = data.NewField(name, nil, make([]*int64, rowIndex+1))
				field.Set(rowIndex, i)
			}

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
	return &Attribute{Name: name, Value: field, TsFormat: datetimeFormat}, nil
}

func (c *Attribute) Size() int {
	return c.Value.Len()
}

func (c *Attribute) Append(value *dynamodb.AttributeValue) error {
	if value.S != nil {
		if c.TsFormat != "" && c.TsFormat != UnixTimestampMiniseconds && c.TsFormat != UnixTimestampSeconds {
			if c.Type() != data.FieldTypeNullableTime {
				return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "S")
			}
			t, err := time.Parse(c.TsFormat, *value.S)
			if err != nil {
				return err
			}
			c.Value.Append(&t)

		} else {
			if c.Type() != data.FieldTypeNullableString {
				return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "S")
			}
			c.Value.Append(value.S)
		}

	} else if value.N != nil {
		i, f, err := parseNumber(*value.N)
		if err != nil {
			return err
		} else if i != nil {
			// int64
			if c.TsFormat == UnixTimestampSeconds {
				if c.Type() != data.FieldTypeNullableTime {
					return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "N")
				}
				t := time.Unix(*i, 0)
				c.Value.Append(&t)
			} else if c.TsFormat == UnixTimestampMiniseconds {
				if c.Type() != data.FieldTypeNullableTime {
					return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "N")
				}

				seconds := *i / 1000
				nanoseconds := (*i % 1000) * 1000000
				t := time.Unix(seconds, nanoseconds)
				c.Value.Append(&t)
			} else if c.TsFormat != "" {
				return errors.New("invalid datetime format")
			} else {
				if c.Type() == data.FieldTypeNullableInt64 {
					c.Value.Append(i)
				} else if c.Type() == data.FieldTypeNullableFloat64 {
					c.Value.Append(aws.Float64(float64(*i)))
				} else {
					return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "N")
				}
			}

		} else {
			// float64
			if c.Type() == data.FieldTypeNullableFloat64 {
				c.Value.Append(f)
			} else if c.Type() == data.FieldTypeNullableInt64 {

				// Convert all previous *int64 values to *float64
				float64Values := make([]*float64, c.Value.Len()+1)
				for i := 0; i < c.Value.Len(); i++ {
					cv, ok := c.Value.ConcreteAt(i)
					if ok {
						float64Values[i] = aws.Float64(float64(cv.(int64)))
					}
				}

				float64Values[c.Value.Len()] = f
				c.Value = data.NewField(c.Name, nil, float64Values)
			} else {
				return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "N")
			}
		}
	} else if value.B != nil {
		if c.Type() != data.FieldTypeNullableString {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "B")
		}
		c.Value.Append(aws.String("[B]"))
	} else if value.BOOL != nil {
		if c.Type() != data.FieldTypeNullableBool {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "BOOL")
		}
		c.Value.Append(value.BOOL)
	} else if value.NULL != nil {
		c.Value.Append(nil)
	} else if value.M != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "M")
		}
		v, err := mapToJson(value)
		if err != nil {
			return err
		}
		c.Value.Append(v)
	} else if value.L != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "L")
		}
		v, err := listToJson(value)
		if err != nil {
			return err
		}
		c.Value.Append(v)
	} else if value.SS != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "SS")
		}
		v, err := stringSetToJson(value)
		if err != nil {
			return err
		}
		c.Value.Append(v)
	} else if value.NS != nil {
		if c.Type() != data.FieldTypeNullableJSON {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "NS")
		}
		v, err := numberSetToJson(value)
		if err != nil {
			return err
		}
		c.Value.Append(v)
	} else if value.BS != nil {
		if c.Type() != data.FieldTypeNullableString {
			return fmt.Errorf("field %s should have type %s, but got %s", c.Name, c.Type().ItemTypeString(), "BS")
		}
		c.Value.Append(aws.String("[BS]"))
	}

	return nil
}
