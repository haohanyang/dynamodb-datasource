import React, { useRef, useState } from "react";
import { Button, CodeEditor, Field, InlineField, InlineFieldRow, Input, Select, TagList } from "@grafana/ui";
import { QueryEditorProps, SelectableValue } from "@grafana/data";
import { DataSource } from "../datasource";
import { DynamoDBDataSourceOptions, DynamoDBQuery, DatetimeFormat } from "../types";
import * as monacoType from "monaco-editor/esm/vs/editor/editor.api";
import "./QueryEditor.css"

type Props = QueryEditorProps<DataSource, DynamoDBQuery, DynamoDBDataSourceOptions>;

const datetimeFormatOptions: Array<SelectableValue<DatetimeFormat>> = [
  {
    label: "ISO 8601",
    value: DatetimeFormat.ISO8601,
    description: "Represents the date and time in UTC, e.g., 2023-05-23T12:34:56Z"
  },
  {
    label: "Unix timestamp",
    value: DatetimeFormat.UnixTimestamp,
    description: "The number of seconds that have elapsed since January 1, 1970 (also known as the Unix epoch), e.g., 1674512096"
  }
]

export function QueryEditor({ query, onChange }: Props) {
  const codeEditorRef = useRef<monacoType.editor.IStandaloneCodeEditor | null>(null);
  const [datetimeFieldInput, setDatetimeFieldInput] = useState<string>("")
  const [datetimeFormatInput, setDatetimeFormatInput] = useState<DatetimeFormat>(DatetimeFormat.ISO8601)

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

  const onAddDatetimeField = (name: string, format: DatetimeFormat) => {
    if (name) {
      onChange({
        ...query,
        datetimeFields: [...query.datetimeFields, { name: name, format: format }]
      });
      setDatetimeFieldInput("")
    }
  }

  const onRemoveDatetimeField = (name: string) => {
    onChange({
      ...query,
      datetimeFields: query.datetimeFields.filter(e => e.name !== name)
    });
  }

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
          <Select options={datetimeFormatOptions} value={datetimeFormatInput} width={30}
            onChange={sv => sv.value && setDatetimeFormatInput(sv.value)}></Select>
        </InlineField>
        <Button onClick={() => onAddDatetimeField(datetimeFieldInput, datetimeFormatInput)}>Add</Button>
      </InlineFieldRow>
      <TagList className="datetime-fields" tags={query.datetimeFields.map(f => f.name)} onClick={(n, _) => onRemoveDatetimeField(n)} />

      <Field label="Query Text" description="The PartiQL statement representing the operation to run">
        <CodeEditor
          onBlur={onQueryTextChange}
          value={query.queryText || ""}
          width="100%"
          height="300px"
          language="sql"
          showLineNumbers={true}
          showMiniMap={false}
          onEditorDidMount={onCodeEditorDidMount}
        />
      </Field>
      <Button onClick={onFormatQueryText}>Format</Button>
    </>
  );
}
