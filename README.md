# Grafana DynamoDB data source

Query your Amazon DynamoDB using PartiQL and visualize the results in your Grafana dashboards.

#### Graph
![screenshot1](/images/screenshot1.png)
#### Table
![screenshot2](/images/screenshot2.png)

## Get started
### Quick start
Run the script [quick_start.py](scripts/quick_start.py) in the root directory to start Grafana containers with the DynamoDB plugin
```
python3 scripts/quick_start.py
```
Visit your Grafana at http://localhost:3000 and configure the data source with your AWS credentials
### Full steps
1. **Download:** Obtain the latest plugin build `haohanyang-dynamodb-datasource-<version>.zip` from the [Release](https://github.com/haohanyang/dynamodb-datasource/releases).

2. **Install:** 
   - Extract the downloaded archive (`haohanyang-dynamodb-datasource-<version>.zip`) into your Grafana plugins directory (`/var/lib/grafana/plugins` or similar).
   - Ensure the plugin binaries (`dynamodb-datasource/gpx_dynamodb_datasource_*`) have execute permissions (`chmod +x`).
### Data source Configuration
The plugin uses [grafana-aws-sdk-react](https://github.com/grafana/grafana-aws-sdk-react) in the configuration page, a common package used for all AWS-related plugins(including plugins made by Grafana Lab). In addition, to test the connection, the plugin requires a "test table", to which the plugin makes a [DescribeTable](https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_DescribeTable.html) request.

### Query data
The plugin currently supports query via [PartiQL](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ql-reference.html). The plugin performs [ExecuteStatement](https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ExecuteStatement.html) on the PartiQL statement that user enters.
#### Datetime attribute
To parse datetime attributes in Grafana, user needs to provide attribute names and format. The format can be unix timestamp (for integers) or [day.js format](https://day.js.org/docs/en/display/format) (for strings)

| Datetime | Format |
| -------- | ------- |
| `1731017392` | Unix timestamp(s) |
| `1731017406839` | Unix timestamp(ms) |
| `2024-10-31T22:04:29+01:00` | `YYYY-MM-DDTHH:mm:ssZ` |
| `2024-10-31T21:04:29Z` | `YYYY-MM-DDTHH:mm:ss[Z]` |
| `2023-08-07T22:18:48.790770` | `YYYY-MM-DDTHH:mm:ss.SSSSSS` |
| `Thu, 31 Oct 2024 21:04:29 GMT` | `ddd, DD MMM YYYY HH:mm:ss z` |

#### Variables
* `$__from` and `$__to` (built-in): start and end in Unix timestamp(ms)
* `$from` and `$to`: start and end in Unix timestamp(s)

You can filter data within the current time range:
```sql
SELECT * FROM MyTable WHERE TimeStamp BETWEEN $from AND $to
```