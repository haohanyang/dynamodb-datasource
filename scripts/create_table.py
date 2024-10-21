import boto3

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
    client.delete_table(TableName=table_name)
    client.create_table(
        TableName=table_name,
        KeySchema=[{"AttributeName": "id", "KeyType": "HASH"}],
        AttributeDefinitions=[
            {
                "AttributeName": "id",
                "AttributeType": "S",
            },
            # {"AttributeName": "myString", "AttributeType": "S"},
            # {"AttributeName": "myNumber", "AttributeType": "N"},
            # {"AttributeName": "myBinary", "AttributeType": "B"},
            # {"AttributeName": "myBool", "AttributeType": "BOOL"},
            # {"AttributeName": "myList", "AttributeType": "L"},
            # {"AttributeName": "myMap", "AttributeType": "M"},
            # {"AttributeName": "myStringSet", "AttributeType": "SS"},
            # {"AttributeName": "myNumberSet", "AttributeType": "NS"},
        ],
        BillingMode="PAY_PER_REQUEST",
    )


def put_item(table_name: str = TABLE_NAME):
    table = get_table(table_name)
    table.put_item(
        Item={
            "id": "1",
            "myString": {"S": "Hello, DynamoDB!"},
            "myNumber": {"N": "123.45"},
            "myBinary": {"B": b"some_binary_data"},
            "myBool": {"BOOL": True},
            "myList": {"L": [{"S": "item1"}, {"N": "2"}, {"BOOL": False}]},
            "myMap": {"M": {"subkey1": {"S": "value1"}, "subkey2": {"N": "99"}}},
            "myStringSet": {"SS": ["value1", "value2", "value3"]},
            "myNumberSet": {"NS": ["1.1", "2.2", "3.3"]},
        }
    )


create_table()
put_item()
