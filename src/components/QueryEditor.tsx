import React, { useRef, useState } from "react";
import { Button, CodeEditor, Field, IconButton, InlineField, InlineFieldRow, Input, Select } from "@grafana/ui";
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
  const [datetimeFormatOption, setDatetimeFormatOption] = useState<string>(DatetimeFormat.UnixTimestampSeconds);
  const [customDatetimeFormatInput, setCustomDatetimeFormatInput] = useState<string>("");

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
    let format = "";
    if (datetimeFormatOption == DatetimeFormat.UnixTimestampMiniseconds || datetimeFormatOption == DatetimeFormat.UnixTimestampSeconds) {
      format = datetimeFormatOption;
    } else {
      format = customDatetimeFormatInput;
    }

    if (datetimeAttributeInput && format && !query.datetimeAttributes.map(e => e.name).includes(datetimeAttributeInput)) {
      onChange({
        ...query,
        datetimeAttributes: [...query.datetimeAttributes, { name: datetimeAttributeInput, format: format }]
      });
      setDatetimeAttributeInput("");
      setCustomDatetimeFormatInput("");
    }
  };

  const showTimeFormat = (format: string) => {
    if (format === DatetimeFormat.UnixTimestampSeconds) {
      return "Unix timestamp(s)";
    } else if (format === DatetimeFormat.UnixTimestampMiniseconds) {
      return "Unix timestamp(ms)";
    } else {
      return format;
    }
  };


  const onRemoveDatetimeAttribute = (name: string) => {
    onChange({
      ...query,
      datetimeAttributes: query.datetimeAttributes.filter(e => e.name !== name)
    });
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
        {datetimeFormatOption === DatetimeFormat.CustomFormat && <InlineField label="Custom Format">
          <Input label="Custom format" placeholder="Enter day.js datatime format" value={customDatetimeFormatInput} onChange={e => setCustomDatetimeFormatInput(e.currentTarget.value)} />
        </InlineField>}
        <Button onClick={onAddDatetimeField} data-testid="datetime-format-add">Add</Button>
      </InlineFieldRow>
      <ul className="datatime-attribute-list">
        {query.datetimeAttributes.map((a, i) =>
          <li className="datatime-attribute-item" key={i}>
            <span className="datatime-attribute-name">{a.name}</span>
            <IconButton name="times" size="lg" tooltip={"Remove \"" + a.name + ": " + showTimeFormat(a.format) + "\""} className="datatime-attribute-remove-btn" onClick={() => onRemoveDatetimeAttribute(a.name)} />
          </li>)}
      </ul>
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
