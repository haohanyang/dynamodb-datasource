package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func LoadExtraPluginSettings(source backend.DataSourceInstanceSettings) (*ExtraPluginSettings, error) {
	settings := ExtraPluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	return &settings, nil
}

func OutputToDataFrame(output *dynamodb.ExecuteStatementOutput) *data.Frame {

	return nil
}
