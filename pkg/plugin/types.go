package plugin

import "github.com/aws/aws-sdk-go/service/dynamodb"

type QueryModel struct {
	QueryText      string
	Limit          int64
	DatetimeFields []DatetimeField
}

type DatetimeField struct {
	Name   string
	Format string
}

var (
	UnixTimestampSeconds     = "unixs"
	UnixTimestampMiniseconds = "unixms"
)

type DynamoDBDataType int

type DataRow map[string]*dynamodb.AttributeValue

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html#HowItWorks.DataTypeDescriptors
const (
	S DynamoDBDataType = iota
	N
	B
	BOOL
	NULL
	M
	L
	SS
	NS
	BS
)

type ExtraPluginSettings struct {
	ConnectionTestTable string `json:"connectionTestTable"`
}
