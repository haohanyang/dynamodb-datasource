import React, { useRef } from "react";
import { Button, CodeEditor, Field, InlineField, Input } from "@grafana/ui";
import { QueryEditorProps } from "@grafana/data";
import { DataSource } from "../datasource";
import { DynamoDBDataSourceOptions, DynamoDBQuery } from "../types";
import * as monacoType from "monaco-editor/esm/vs/editor/editor.api";

type Props = QueryEditorProps<DataSource, DynamoDBQuery, DynamoDBDataSourceOptions>;

export function QueryEditor({ query, onChange }: Props) {
  const codeEditorRef = useRef<monacoType.editor.IStandaloneCodeEditor | null>(null);

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

  const onCodeEditorDidMount = (e: monacoType.editor.IStandaloneCodeEditor) => {
    codeEditorRef.current = e;
  };

  return (
    <>
      <InlineField label="Limit" tooltip="(Optional) The maximum number of items to evaluate">
        <Input type="number" min={0} value={query.limit} onChange={onLimitChange} />
      </InlineField>
      <Field label="Query Text" description="The PartiQL statement representing the operation to run">
        <CodeEditor
          onChange={onQueryTextChange}
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
