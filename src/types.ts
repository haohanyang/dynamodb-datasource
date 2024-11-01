import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from "@grafana/aws-sdk";
import { DataQuery } from "@grafana/schema";

export interface DynamoDBQuery extends DataQuery {
  queryText?: string;
  limit?: number;
  datetimeAttributes: DatetimeAttribute[];
}

export const DEFAULT_QUERY: Partial<DynamoDBQuery> = {
  queryText: "",
  datetimeAttributes: []
};

export interface DynamoDBDataSourceOptions extends AwsAuthDataSourceJsonData {
  connectionTestTable?: string;
}

export interface DynamoDBDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData { }

export const DatetimeFormat = {
  UnixTimestampSeconds: "1",
  UnixTimestampMiniseconds: "2",
  CustomFormat: "custom"
};
export interface DatetimeAttribute {
  name: string;
  format: string;
}
