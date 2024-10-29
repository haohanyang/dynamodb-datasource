import React, { useRef, useState } from "react";
import { Button, CodeEditor, Field, InlineField, InlineFieldRow, Input, Select, TagList } from "@grafana/ui";
import { QueryEditorProps, SelectableValue } from "@grafana/data";
import { DataSource } from "../datasource";
import { DynamoDBDataSourceOptions, DynamoDBQuery, DatetimeFormat } from "../types";
import * as monacoType from "monaco-editor/esm/vs/editor/editor.api";
import "./QueryEditor.css";

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
    label: "RFC 3339",
    value: DatetimeFormat.RFC3339,
    description: DatetimeFormat.RFC3339
  },
  {
    label: "RFC 3339 nano",
    value: DatetimeFormat.RFC3339Nano,
    description: DatetimeFormat.RFC3339Nano
  },
  {
    label: "RFC 1123",
    value: DatetimeFormat.RFC1123,
    description: DatetimeFormat.RFC1123
  },
  {
    label: "RFC 1123Z",
    value: DatetimeFormat.RFC1123Z,
    description: DatetimeFormat.RFC1123Z
  },
  {
    label: "RFC 822",
    value: DatetimeFormat.RFC822,
    description: DatetimeFormat.RFC822
  },
  {
    label: "RFC 822Z",
    value: DatetimeFormat.RFC822Z,
    description: DatetimeFormat.RFC822Z
  },
  {
    label: "RFC 850",
    value: DatetimeFormat.RFC850,
    description: DatetimeFormat.RFC850
  },
  {
    label: "Custom format",
    value: DatetimeFormat.CustomFormat,
    description: "User-defined format"
  }
];

// const formatHelper = "User-defined format based on the specific timestamp (January 2, 15:04:05, 2006, in time zone seven hours west of GMT), e.g. \"2006-01-02T15:04:05Z07:00\"(RFC3339), \"Mon, 02 Jan 2006 15:04:05 MST\"(RFC1123), \"02 Jan 06 15:04 MST\"(RFC822)"

export function QueryEditor({ query, onChange }: Props) {
  const codeEditorRef = useRef<monacoType.editor.IStandaloneCodeEditor | null>(null);
  const [datetimeFieldInput, setDatetimeFieldInput] = useState<string>("");
  const [datetimeFormatOption, setDatetimeFormatOption] = useState<string>(DatetimeFormat.UnixTimestampMiniseconds);
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

  const onAddDatetimeField = (name: string, option: string, customFormat: string) => {
    let format = customFormat;
    if (option !== DatetimeFormat.CustomFormat) {
      format = option;
    }

    if (name && format) {
      onChange({
        ...query,
        datetimeFields: [...query.datetimeFields, { name: name, format: format }]
      });
      setDatetimeFieldInput("");
      setCustomDatetimeFormatInput("");
    }
  };

  const onRemoveDatetimeField = (name: string) => {
    onChange({
      ...query,
      datetimeFields: query.datetimeFields.filter(e => e.name !== name)
    });
  };

  const onCodeEditorDidMount = (e: monacoType.editor.IStandaloneCodeEditor) => {
    codeEditorRef.current = e;
  };

  return (
    <>
      <InlineField label="Limit" tooltip="(Optional) The maximum number of items to evaluate">
        <Input type="number" min={0} value={query.limit} onChange={onLimitChange} aria-label="Limit" />
      </InlineField>
      <InlineFieldRow label="Add datetime field">
        <InlineField label="Field" tooltip="Field which has datetime data type">
          <Input value={datetimeFieldInput} onChange={e => setDatetimeFieldInput(e.currentTarget.value)} />
        </InlineField>
        <InlineField label="Format">
          <Select options={datetimeFormatOptions} value={datetimeFormatOption} width={30}
            onChange={sv => sv.value && setDatetimeFormatOption(sv.value)}></Select>
        </InlineField>
        {datetimeFormatOption === DatetimeFormat.CustomFormat && <InlineField label="Custom Format">
          <Input label="Custom format" value={customDatetimeFormatInput} onChange={e => setCustomDatetimeFormatInput(e.currentTarget.value)} />
        </InlineField>}
        <Button onClick={() => onAddDatetimeField(datetimeFieldInput, datetimeFormatOption, customDatetimeFormatInput)}>Add</Button>
      </InlineFieldRow>
      <TagList className="datetime-fields" tags={query.datetimeFields.map(f => f.name)} onClick={(n, _) => onRemoveDatetimeField(n)} />
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
