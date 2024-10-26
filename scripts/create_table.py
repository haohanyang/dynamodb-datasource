import time
import boto3
from datetime import datetime
from decimal import Decimal
from pprint import pprint

TABLE_NAME = "test"


def get_client():
    client = boto3.client(
        "dynamodb",
        endpoint_url="http://localhost:4566",
        aws_access_key_id="test",
        region_name="us-east-1",
        aws_secret_access_key="test",
    )
    return client


def get_table(table_name: str):
    resource = boto3.resource(
        "dynamodb",
        endpoint_url="http://localhost:4566",
        aws_access_key_id="test",
        region_name="us-east-1",
        aws_secret_access_key="test",
    )
    return resource.Table(table_name)


def create_table(table_name: str = TABLE_NAME):
    client = get_client()

    tables = client.list_tables()["TableNames"]
    if table_name in tables:
        client.delete_table(TableName=table_name)

    client.create_table(
        TableName=table_name,
        KeySchema=[
            {"AttributeName": "id", "KeyType": "HASH"},
            {"AttributeName": "sid", "KeyType": "RANGE"},
        ],
        AttributeDefinitions=[
            {
                "AttributeName": "id",
                "AttributeType": "N",
            },
            {
                "AttributeName": "sid",
                "AttributeType": "S",
            },
        ],
        BillingMode="PAY_PER_REQUEST",
    )


def put_item(table_name: str = TABLE_NAME):
    table = get_table(table_name)
    table.put_item(
        Item={
            "id": 1,
            "sid": "A",
            "myString": "Hello, DynamoDB!",
            "myFloat": Decimal("123.45"),
            "myInt": 123,
            "myBinary": b"some_binary_data",
            "myBool": True,
            "myList": ["item1", 2, False],
            "myMap": {"subkey1": "value1", "subkey2": 99},
            "myStringSet": set(["value1", "value2", "value3"]),
            "myNumberSet": set([Decimal("1.1"), 2, Decimal("3.3")]),
            "myISODate": datetime.now().isoformat(),
            "myUnixDate": Decimal(str(time.mktime(datetime.now().timetuple()))),
        }
    )
    item = table.get_item(Key={"id": 1, "sid": "A"})["Item"]
    pprint(item)


create_table()
put_item()
