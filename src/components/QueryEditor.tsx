import React, { useRef } from "react";
import { Button, CodeEditor, Field } from "@grafana/ui";
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
      <Field label="Query Text" description="PartiQL text">
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
