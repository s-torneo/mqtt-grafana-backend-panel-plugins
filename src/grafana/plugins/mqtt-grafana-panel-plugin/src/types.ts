import { DataSourceJsonData } from '@grafana/data';

export interface SimpleOptions {
  mqttTopic: string;
  color_text: string;
  color_button: string;
  datasource: number;
  operation: string;
  publishMsg: string;
  buttonName: string;
  backgroundcolor_table: string;
  textcolor_table: string;
}

export interface DataSourceOptions extends DataSourceJsonData {}
