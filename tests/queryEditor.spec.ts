import { CreateTableCommand, DeleteTableCommand, DynamoDBClient, ListTablesCommand, PutItemCommand } from '@aws-sdk/client-dynamodb';
import { test, expect } from '@grafana/plugin-e2e';

test.beforeAll(async function ({ createDataSource, readProvisionedDataSource }) {
    process.env.AWS_ACCESS_KEY_ID = "test"
    process.env.AWS_SECRET_ACCESS_KEY = "test"
    const client = new DynamoDBClient({
        endpoint: "http://localhost:4566",
        region: "us-east-1",
    });

    const res = await client.send(new ListTablesCommand())
    if (res.TableNames && res.TableNames.includes("test")) {
        await client.send(new DeleteTableCommand({
            TableName: "test"
        }))
    }

    await client.send(new CreateTableCommand({
        TableName: "test",
        AttributeDefinitions: [
            {
                AttributeName: "id",
                AttributeType: "N",
            },
        ],
        KeySchema: [
            {
                AttributeName: "id",
                KeyType: "HASH",
            },
        ],
        ProvisionedThroughput: {
            ReadCapacityUnits: 1,
            WriteCapacityUnits: 1,
        },
    }))

    await client.send(new PutItemCommand({
        TableName: "test",
        Item: {
            id: { N: "1" },
        }
    }))

    await client.send(new PutItemCommand({
        TableName: "test",
        Item: {
            id: { N: "2" },
        }
    }))

    await client.send(new PutItemCommand({
        TableName: "test",
        Item: {
            id: { N: "3" },
        }
    }))

    const ds = await readProvisionedDataSource({ fileName: "e2e.yml" });
    await createDataSource(ds)
})

test("should return correct query result", async ({
    panelEditPage,
    readProvisionedDataSource,
    selectors
}) => {
    const ds = await readProvisionedDataSource({ fileName: "e2e.yml" });
    await panelEditPage.datasource.set(ds.name);
    await panelEditPage.getQueryEditorRow("A").getByLabel("Limit").fill("3");

    const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container, {
        root: panelEditPage.getQueryEditorRow("A")
    }).getByRole("textbox");

    await editor.clear();
    await editor.fill("SELECT * FROM test");

    await panelEditPage.setVisualization('Table');
    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.data).toContainText(["2", "1", "3"]);
});
