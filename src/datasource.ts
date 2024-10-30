import { DataSourceInstanceSettings, CoreApp, ScopedVars, DataQueryRequest, DataQueryResponse } from "@grafana/data";
import { DataSourceWithBackend, getTemplateSrv } from "@grafana/runtime";
import { Observable } from "rxjs";
import { DynamoDBQuery, DynamoDBDataSourceOptions, DEFAULT_QUERY, DatetimeFormat } from "./types";
import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
import utc from "dayjs/plugin/utc";

dayjs.extend(utc);
dayjs.extend(timezone);

export class DataSource extends DataSourceWithBackend<DynamoDBQuery, DynamoDBDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<DynamoDBDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<DynamoDBQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: DynamoDBQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      queryText: getTemplateSrv().replace(query.queryText, scopedVars),
    };
  }

  filterQuery(query: DynamoDBQuery): boolean {
    // if no query has been provided, prevent the query from being executed
    return !!query.queryText;
  }

  query(request: DataQueryRequest<DynamoDBQuery>): Observable<DataQueryResponse> {
    // Golang's reference time: Mon, 02 Jan 2006 15:04:05 MST
    const refTime = dayjs.tz(1136239445999, "America/Edmonton");
    const queries = request.targets.map((query) => {
      return {
        ...query,
        queryText:
          query.queryText?.replaceAll(/\$from/g, Math.floor(request.range.from.toDate().getTime() / 1000).toString())
            .replaceAll(/\$to/g, Math.floor(request.range.to.toDate().getTime() / 1000).toString()),
        datetimeFields: query.datetimeFields.map(field => {
          if (field.format != DatetimeFormat.UnixTimestampSeconds && field.format != DatetimeFormat.UnixTimestampMiniseconds) {
            return { ...field, format: refTime.format(field.format) };
          }
          return field;
        })
      };
    });
    return super.query({ ...request, targets: queries });
  }
}
