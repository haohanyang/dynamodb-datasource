import time
import boto3
import random
import time
import sys
from datetime import datetime, timedelta
from decimal import Decimal


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


def create_table(table_name: str):
    client = get_client()

    tables = client.list_tables()["TableNames"]
    if table_name in tables:
        client.delete_table(TableName=table_name)

    client.create_table(
        TableName=table_name,
        KeySchema=[
            {"AttributeName": "Id", "KeyType": "HASH"},
        ],
        AttributeDefinitions=[
            {
                "AttributeName": "Id",
                "AttributeType": "N",
            }
        ],
        BillingMode="PAY_PER_REQUEST",
    )


def create_general_table(table_name: str = "Test"):
    create_table(table_name)

    table = get_table(table_name)
    table.put_item(
        Item={
            "Id": 1,
            "String": "Hello, DynamoDB!",
            "Float": Decimal("123.45"),
            "Int": 123,
            "Binary": b"some_binary_data",
            "Bool": True,
            "List": ["item1", 2, False],
            "Map": {"subkey1": "value1", "subkey2": 99},
            "StringSet": set(["value1", "value2", "value3"]),
            "NumberSet": set([Decimal("1.1"), 2, Decimal("3.3")]),
            "ISODate": datetime.now().isoformat(),
            "UnixDate": Decimal(str(time.mktime(datetime.now().timetuple()))),
        }
    )


def create_ts_table(table_name: str = "TestTimeSeries"):
    create_table(table_name)

    # generate random ts data
    data = []
    current_ts = datetime.now()
    current_vals = [0.0, 0.0]
    for i in range(10000):
        data.append((current_ts, "A", current_vals[0]))
        data.append((current_ts, "B", current_vals[1]))
        current_vals[0] += random.randint(-10, 10) * random.random()
        current_vals[1] += random.random()
        current_ts -= timedelta(minutes=1)

    table = get_table(table_name)

    with table.batch_writer() as batch:
        for i in range(10000):
            ts: datetime = data[i][0]
            value: float = data[i][2]
            batch.put_item(
                Item={
                    "Id": i,
                    "Time": int(time.mktime(ts.timetuple())),
                    "Type": data[i][1],
                    "Val": Decimal(str(value)),
                }
            )


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Invalid input")
        exit(1)

    if sys.argv[1] == "table":
        if len(sys.argv) > 2:
            create_general_table(sys.argv[2])
        else:
            create_general_table()
    elif sys.argv[1] == "ts":
        if len(sys.argv) > 2:
            create_ts_table(sys.argv[2])
        else:
            create_ts_table()
    else:
        print("Invalid input")
        exit(1)
