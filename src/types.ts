import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from "@grafana/aws-sdk";
import { DataQuery } from "@grafana/schema";

export interface DynamoDBQuery extends DataQuery {
  queryText?: string;
  limit?: number
  datetimeFields: DatetimeField[]
}

export const DEFAULT_QUERY: Partial<DynamoDBQuery> = {
  queryText: "",
  datetimeFields: []
};

export interface DynamoDBDataSourceOptions extends AwsAuthDataSourceJsonData {
  connectionTestTable?: string;
}

export interface DynamoDBDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData { }

export enum DatetimeFormat {
  UnixTimestampSeconds = 1,
  UnixTimestampMiniseconds = 2,
  CustomFormat = 3
}
export interface DatetimeField {
  name: string
  format: string
}
