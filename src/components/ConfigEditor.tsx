import React from "react";
import { ConnectionConfig } from "@grafana/aws-sdk";
import { DataSourcePluginOptionsEditorProps } from "@grafana/data";
import { DynamoDBDataSourceOptions, DynamoDBDataSourceSecureJsonData } from "../types";
import { Field, Input } from "@grafana/ui";

interface Props extends DataSourcePluginOptionsEditorProps<DynamoDBDataSourceOptions, DynamoDBDataSourceSecureJsonData> { }

const standardRegions = [
  "af-south-1",
  "ap-east-1",
  "ap-northeast-1",
  "ap-northeast-2",
  "ap-northeast-3",
  "ap-south-1",
  "ap-southeast-1",
  "ap-southeast-2",
  "ap-southeast-3",
  "ca-central-1",
  "cn-north-1",
  "cn-northwest-1",
  "eu-central-1",
  "eu-north-1",
  "eu-south-1",
  "eu-west-1",
  "eu-west-2",
  "eu-west-3",
  "me-south-1",
  "sa-east-1",
  "us-east-1",
  "us-east-2",
  "us-gov-east-1",
  "us-gov-west-1",
  "us-west-1",
  "us-west-2",
];

export function ConfigEditor(props: Props) {
  const onTestTableChange: React.FormEventHandler<HTMLInputElement> = e => {
    props.onOptionsChange({
      ...props.options,
      jsonData:
      {
        ...props.options.jsonData,
        connectionTestTable: e.currentTarget.value
      }
    });
  };

  return (
    <div className="width-30">
      <ConnectionConfig {...props} standardRegions={standardRegions} />
      <Field label="Test table" description="Name of table for connection test">
        <Input value={props.options.jsonData.connectionTestTable} onChange={onTestTableChange} aria-label="Test table"></Input>
      </Field>
    </div>
  );
};
