import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';
import { DataQuery } from '@grafana/schema';

export interface DynamoDBQuery extends DataQuery {
  queryText?: string;
}

export const DEFAULT_QUERY: Partial<DynamoDBQuery> = {
  queryText: "",
};

export interface DynamoDBDataSourceOptions extends AwsAuthDataSourceJsonData {
  connectionTestTable?: string;
}

export interface DynamoDBDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData { }


