package plugin

type QueryModel struct {
	QueryText string
}

type DynamoDBDataType int

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
