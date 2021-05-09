import { DataSourceJsonData } from '@grafana/data';

export interface SimpleOptions {
  topic: string;
  color_text: string;
  color_button: string;
  datasource: number;
  operation: string;
}

export type GetDataResponse = {
  payload: string;
};

export interface DataSourceOptions extends DataSourceJsonData {}
