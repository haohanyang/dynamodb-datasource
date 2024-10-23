import { test, expect } from "@grafana/plugin-e2e";
import { DynamoDBDataSourceOptions, DynamoDBDataSourceSecureJsonData } from "../src/types";
import { CreateTableCommand, DeleteTableCommand, DynamoDBClient, ListTablesCommand } from "@aws-sdk/client-dynamodb";

test.beforeAll(async function () {
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
})

test(`"Save & test" should be successful when configuration is valid`, async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<DynamoDBDataSourceOptions, DynamoDBDataSourceSecureJsonData>({ fileName: 'e2e.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });

  await page.getByLabel("Authentication Provider", { exact: true }).fill("Access & secret key")
  await page.keyboard.press("Enter")
  await page.getByLabel("Access Key ID").fill("test")
  await page.getByLabel("Secret Access Key").fill("test")
  await page.getByLabel("Endpoint").fill("http://localstack:4566")
  await page.getByLabel("Default Region").fill("us-east-1")
  await page.keyboard.press("Enter")
  await page.getByLabel("Test table").fill("test")
  await expect(configPage.saveAndTest()).toBeOK();
});

test(`"Save & test" should be successful when test table doesn't exist`, async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<DynamoDBDataSourceOptions, DynamoDBDataSourceSecureJsonData>({ fileName: 'e2e.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });

  await page.getByLabel("Authentication Provider", { exact: true }).fill("Access & secret key")
  await page.keyboard.press("Enter")
  await page.getByLabel("Access Key ID").fill("test")
  await page.getByLabel("Secret Access Key").fill("test")
  await page.getByLabel("Endpoint").fill("http://localstack:4566")
  await page.getByLabel("Default Region").fill("us-east-1")
  await page.keyboard.press("Enter")
  await page.getByLabel("Test table").fill("null")
  await expect(configPage.saveAndTest()).not.toBeOK();
});
