package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	dsSetting := awsds.AWSDatasourceSettings{}
	err := dsSetting.Load(settings)
	if err != nil {
		backend.Logger.Error("failed to load settings", err.Error())
		return nil, err
	}

	authSettings := awsds.ReadAuthSettings(ctx)
	sessionCache := awsds.NewSessionCache()

	return &Datasource{
		Settings:     dsSetting,
		authSettings: *authSettings,
		sessionCache: sessionCache,
	}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	Settings     awsds.AWSDatasourceSettings
	sessionCache *awsds.SessionCache
	authSettings awsds.AuthSettings
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

func (d *Datasource) getDynamoDBClient(ctx context.Context, settings *backend.DataSourceInstanceSettings) (*dynamodb.DynamoDB, error) {
	httpClientProvider := httpclient.NewProvider()
	httpClientOptions, err := settings.HTTPClientOptions(ctx)
	if err != nil {
		backend.Logger.Error("failed to create http client options", err.Error())
		return nil, err
	}

	httpClient, err := httpClientProvider.New(httpClientOptions)
	if err != nil {
		backend.Logger.Error("failed to create http client", err.Error())
		return nil, err
	}

	session, err := d.sessionCache.GetSessionWithAuthSettings(awsds.GetSessionConfig{
		Settings:      d.Settings,
		HTTPClient:    httpClient,
		UserAgentName: aws.String("DynamoDB"),
	}, d.authSettings)

	if err != nil {
		return nil, err
	}

	return dynamodb.New(session), nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	dynamoDBClient, err := d.getDynamoDBClient(ctx, req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		return nil, err
	}

	for _, q := range req.Queries {
		res := d.query(ctx, dynamoDBClient, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *Datasource) query(ctx context.Context, dynamoDBClient *dynamodb.DynamoDB, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	var qm QueryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(qm.QueryText),
	}

	if qm.Limit > 0 {
		input.Limit = aws.Int64(qm.Limit)
	}

	output, err := dynamoDBClient.ExecuteStatementWithContext(ctx, input)

	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("executes statement: %v", err.Error()))
	}

	dateFields := make(map[string]string)
	for _, k := range qm.DatetimeAttributes {
		dateFields[k.Name] = k.Format
	}

	frame, err := QueryResultToDataFrame(query.RefID, output, dateFields)
	if err != nil {
		response.Error = err
		return response
	}

	response.Frames = append(response.Frames, frame)
	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	backend.Logger.Debug("Checking health")
	res := &backend.CheckHealthResult{}

	client, err := d.getDynamoDBClient(ctx, req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = err.Error()
		return res, nil
	}

	extraSettings, err := loadExtraPluginSettings(*req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = err.Error()
		return res, err
	}

	_, err = client.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(extraSettings.ConnectionTestTable),
	})

	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = err.Error()
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Successfully connects to DynamoDB",
	}, nil
}
