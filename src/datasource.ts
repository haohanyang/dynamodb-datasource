import { DataSourceInstanceSettings, CoreApp, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { DynamoDBQuery, DynamoDBDataSourceOptions, DEFAULT_QUERY } from './types';

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
}
