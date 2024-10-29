import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from "@grafana/aws-sdk";
import { DataQuery } from "@grafana/schema";

export interface DynamoDBQuery extends DataQuery {
  queryText?: string;
  limit?: number;
  datetimeFields: DatetimeField[];
}

export const DEFAULT_QUERY: Partial<DynamoDBQuery> = {
  queryText: "",
  datetimeFields: []
};

export interface DynamoDBDataSourceOptions extends AwsAuthDataSourceJsonData {
  connectionTestTable?: string;
}

export interface DynamoDBDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData { }

export const DatetimeFormat = {
  UnixTimestampSeconds: "1",
  UnixTimestampMiniseconds: "2",
  RFC822: "02 Jan 06 15:04 MST",
  RFC822Z: "02 Jan 06 15:04 -0700",
  RFC850: "Monday, 02-Jan-06 15:04:05 MST",
  RFC1123: "Mon, 02 Jan 2006 15:04:05 MST",
  RFC1123Z: "Mon, 02 Jan 2006 15:04:05 -0700",
  RFC3339: "2006-01-02T15:04:05Z07:00",
  RFC3339Nano: "2006-01-02T15:04:05.999999999Z07:00",
  CustomFormat: "custom"
};
export interface DatetimeField {
  name: string;
  format: string;
}
