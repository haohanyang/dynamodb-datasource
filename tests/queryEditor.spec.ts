import { AttributeValue, CreateTableCommand, DeleteTableCommand, DynamoDBClient, ListTablesCommand, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { PanelEditPage, test, expect } from "@grafana/plugin-e2e";
import { DatetimeAttribute, DatetimeFormat } from "../src/types";
import { Page } from "@playwright/test";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

function getDynamoClient() {
    return new DynamoDBClient({
        endpoint: "http://localhost:4566",
        region: "us-east-1",
        credentials: {
            accessKeyId: "test",
            secretAccessKey: "test"
        }
    });
}

async function initTableWithItems(tableName: string, items: Array<Record<string, AttributeValue>>) {
    const client = getDynamoClient();

    const res = await client.send(new ListTablesCommand());
    if (res.TableNames && res.TableNames.includes(tableName)) {
        await client.send(new DeleteTableCommand({
            TableName: tableName
        }));
    }

    await client.send(new CreateTableCommand({
        TableName: tableName,
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
            }
        ],
        BillingMode: "PAY_PER_REQUEST"
    }));

    items.map(async (item, idx) => {
        await client.send(new PutItemCommand({
            TableName: tableName,
            Item: {
                id: { N: idx.toString() },
                ...item
            }
        }));
    });
}

async function addDatetimeFormats(datetimeAttributes: DatetimeAttribute[], panelEditPage: PanelEditPage, page: Page) {
    for (const attribute of datetimeAttributes) {
        // Add attribute name
        await panelEditPage.getQueryEditorRow("A").getByTestId("datetime-attribute-input").fill(attribute.name);

        if (attribute.format === DatetimeFormat.UnixTimestampMiniseconds || attribute.format === DatetimeFormat.UnixTimestampSeconds) {

            if (attribute.format === DatetimeFormat.UnixTimestampMiniseconds) {
                await panelEditPage.getQueryEditorRow("A").getByTestId("datetime-format-select-input").fill("Unix timestamp(ms)");
            } else {
                await panelEditPage.getQueryEditorRow("A").getByTestId("datetime-format-select-input").fill("Unix timestamp(s)");
            }
            await page.keyboard.press("Enter");
            await panelEditPage.getQueryEditorRow("A").getByTestId("datetime-format-add").click();
        } else {
            await panelEditPage.getQueryEditorRow("A").getByTestId("datetime-format-select-input").fill("Custom format");
            await page.keyboard.press("Enter");
            await new Promise((resolve) => {
                setTimeout(resolve, 1000);
            });
            await panelEditPage.getQueryEditorRow("A").getByPlaceholder("Enter day.js datatime format").fill(attribute.format);
            await panelEditPage.getQueryEditorRow("A").getByTestId("datetime-format-add").click();
        }
    }
}

test.beforeAll(async function ({ createDataSource, readProvisionedDataSource }) {
    const ds = await readProvisionedDataSource({ fileName: "e2e.yml" });
    await createDataSource(ds);
});


test("should return correct datetime", async ({
    panelEditPage,
    readProvisionedDataSource,
    selectors,
    page
}) => {
    test.setTimeout(100000)
    await initTableWithItems("test", [
        {
            "t1": { N: "1730408669" }, "t2": { S: "2024-10-31T22:04:29+01:00" },
            "t3": { S: "Thu, 31 Oct 2024 21:04:29 GMT" }, "t4": { S: "2024-10-31T21:04:29Z" }
        },
    ]);
    const ds = await readProvisionedDataSource({ fileName: "e2e.yml" });
    await panelEditPage.datasource.set(ds.name);

    await panelEditPage.setVisualization("Table");

    await addDatetimeFormats([
        { name: "t1", format: DatetimeFormat.UnixTimestampSeconds },
        { name: "t2", format: "YYYY-MM-DDTHH:mm:ssZ" },
        { name: "t3", format: "ddd, DD MMM YYYY HH:mm:ss z" },
        { name: "t4", format: "YYYY-MM-DDTHH:mm:ss[Z]" }
    ], panelEditPage, page);

    const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container, {
        root: panelEditPage.getQueryEditorRow("A")
    }).getByRole("textbox");

    await editor.clear();
    await editor.fill("SELECT * FROM test");

    await expect(panelEditPage.refreshPanel()).toBeOK();

    const expected = dayjs.tz(1730408669000, Intl.DateTimeFormat().resolvedOptions().timeZone).format("YYYY-MM-DD HH:mm:ss");
    console.log(expected);
    await expect(panelEditPage.panel.data.getByText(expected)).toHaveCount(4);
});
