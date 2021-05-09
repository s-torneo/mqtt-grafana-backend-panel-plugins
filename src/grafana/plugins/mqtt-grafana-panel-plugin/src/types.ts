import { DataSourceJsonData } from '@grafana/data';

export interface SimpleOptions {
  mqttTopic: string;
  color_text: string;
  color_button: string;
  datasource: number;
  operation: string;
  publishMsg: string;
}

export type GetDataResponse = {
  payload: string;
};

export interface DataSourceOptions extends DataSourceJsonData {}
