import React, { useRef, useState } from "react";
import { Button, CodeEditor, Field, InlineField, InlineFieldRow, Input, Select, TagsInput } from "@grafana/ui";
import { QueryEditorProps, SelectableValue } from "@grafana/data";
import { DataSource } from "../datasource";
import { DynamoDBDataSourceOptions, DynamoDBQuery, DatetimeFormat } from "../types";
import * as monacoType from "monaco-editor/esm/vs/editor/editor.api";
import "./QueryEditor.css";
import { Divider } from "@grafana/aws-sdk";

type Props = QueryEditorProps<DataSource, DynamoDBQuery, DynamoDBDataSourceOptions>;

const datetimeFormatOptions: Array<SelectableValue<string>> = [
  {
    label: "Unix timestamp(s)",
    value: DatetimeFormat.UnixTimestampSeconds,
    description: "The number of seconds that have elapsed since January 1, 1970 (also known as the Unix epoch), e.g., 1674512096"
  },
  {
    label: "Unix timestamp(ms)",
    value: DatetimeFormat.UnixTimestampMiniseconds,
    description: "The number of miliseconds that have elapsed since January 1, 1970 (also known as the Unix epoch), e.g., 1674512096000"
  },
  {
    label: "Custom format",
    value: DatetimeFormat.CustomFormat,
    description: "User-defined format (moment.js/day.js)"
  }
];


export function QueryEditor({ query, onChange }: Props) {
  const codeEditorRef = useRef<monacoType.editor.IStandaloneCodeEditor | null>(null);
  const [datetimeAttributeInput, setDatetimeAttributeInput] = useState<string>("");
  const [datetimeFormatOption, setDatetimeFormatOption] = useState<string>(DatetimeFormat.UnixTimestampMiniseconds);

  const onQueryTextChange = (text: string) => {
    onChange({ ...query, queryText: text });
  };

  const onLimitChange: React.FormEventHandler<HTMLInputElement> = e => {
    if (!e.currentTarget.value) {
      onChange({ ...query, limit: undefined });
    } else {
      const parsed = Number.parseInt(e.currentTarget.value, 10);
      if (Number.isInteger(parsed) && parsed > 0) {
        onChange({ ...query, limit: parsed });
      }
    }
  };


  const onFormatQueryText = () => {
    if (codeEditorRef.current) {
      codeEditorRef.current.getAction("editor.action.formatDocument").run();
    }
  };

  const onAddDatetimeField = () => {
    if (datetimeAttributeInput && (datetimeFormatOption === DatetimeFormat.UnixTimestampMiniseconds || datetimeFormatOption === DatetimeFormat.UnixTimestampSeconds)
      && !query.datetimeAttributes.map(e => e.name).includes(datetimeAttributeInput)) {
      onChange({
        ...query,
        datetimeAttributes: [...query.datetimeAttributes, { name: datetimeAttributeInput, format: datetimeFormatOption }]
      });
      setDatetimeAttributeInput("");
    }
  };


  const onDatetimeFormatsChange = (formats: string[]) => {
    if (formats.length < query.datetimeAttributes.length) {
      // Remove an attribute
      onChange({
        ...query,
        datetimeAttributes: query.datetimeAttributes.filter(e => formats.includes(e.name))
      });
    } else if (formats.length > query.datetimeAttributes.length && datetimeAttributeInput) {
      // Add an attribute
      const format = formats[formats.length - 1]; //formats.find(e => !query.datetimeAttributes.map(f => f.name).includes(e));
      if (!query.datetimeAttributes.map(e => e.name).includes(datetimeAttributeInput)) {
        onChange({
          ...query,
          datetimeAttributes: [...query.datetimeAttributes, { name: datetimeAttributeInput, format: format }]
        });
        setDatetimeAttributeInput("");
      }
    }


  };

  const onCodeEditorDidMount = (e: monacoType.editor.IStandaloneCodeEditor) => {
    codeEditorRef.current = e;
  };

  return (
    <>
      <InlineFieldRow>
        <InlineField label="Limit" tooltip="(Optional) The maximum number of items to evaluate" labelWidth={11}>
          <Input type="number" min={0} value={query.limit} onChange={onLimitChange} aria-label="Limit" width={15} />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Attribute" tooltip="Attribute which has datetime data type" labelWidth={11}>
          <Input value={datetimeAttributeInput} onChange={e => setDatetimeAttributeInput(e.currentTarget.value)}
            data-testid="datetime-attribute-input" width={15} />
        </InlineField>
        <InlineField label="Format" labelWidth={11}>
          <Select options={datetimeFormatOptions} value={datetimeFormatOption} width={25}
            onChange={sv => sv.value && setDatetimeFormatOption(sv.value)} data-testid="datetime-format-select" />
        </InlineField>
        {datetimeFormatOption !== DatetimeFormat.CustomFormat && <Button onClick={onAddDatetimeField} data-testid="datetime-format-add">Add2</Button>}
      </InlineFieldRow>
      <TagsInput tags={query.datetimeAttributes.map(f => f.name)} onChange={onDatetimeFormatsChange}
        disabled={datetimeFormatOption !== DatetimeFormat.CustomFormat} width={38} placeholder="Enter day.js datatime format" />
      <Divider />
      <Field label="Query Text" description="The PartiQL statement representing the operation to run">
        <CodeEditor
          onBlur={onQueryTextChange}
          value={query.queryText || ""}
          width="100%"
          height="100px"
          language="sql"
          showMiniMap={false}
          monacoOptions={{ fontSize: 14 }}
          onEditorDidMount={onCodeEditorDidMount}
        />
      </Field>
      <Button onClick={onFormatQueryText}>Format</Button>
    </>
  );
}
