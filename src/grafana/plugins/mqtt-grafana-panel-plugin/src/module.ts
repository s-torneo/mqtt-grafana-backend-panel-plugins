import { PanelPlugin } from '@grafana/data';
import { SimpleOptions } from './types';
import { MqttPanel } from './MqttPanel';
import * as c from './constants';

export const plugin = new PanelPlugin<SimpleOptions>(MqttPanel).useFieldConfig().setPanelOptions((builder) => {
  return builder
    .addSelect({
      path: 'operation',
      name: 'MQTT Operation',
      description: 'Select operation',
      defaultValue: 'connect',
      settings: {
        options: [
          { label: c.connectName, value: c.connectOp },
          { label: c.disconnectName, value: c.disconnectOp },
          { label: c.subName, value: c.subOp },
          { label: c.unsubscribeName, value: c.unsubscribeOp },
          { label: c.publishName, value: c.publishOp },
          { label: c.getDataName, value: c.getDataOp },
        ],
      },
    })
    .addTextInput({
      path: 'mqttTopic',
      name: 'Topic',
      description: 'Insert MQTT Topic (Grafana variables supported)',
      defaultValue: 'test',
    })
    .addTextInput({
      path: 'publishMsg',
      name: 'Message',
      description: 'Insert MQTT Message',
      defaultValue: '',
    })
    .addTextInput({
      path: 'buttonName',
      name: 'Button Name',
      description: 'Insert the button name',
      defaultValue: '',
    })
    .addColorPicker({
      path: 'color_button',
      name: 'Button Color',
      description: 'Choose the color of the button',
      defaultValue: '',
    })
    .addColorPicker({
      path: 'color_text',
      name: 'Button Text Color',
      description: 'Choose the text color of the button',
      defaultValue: '',
    })
    .addTextInput({
      path: 'datasource',
      name: 'Datasource Id',
      description: 'Id of the datasource (corresponds to the mqtt client configured)',
      defaultValue: '1',
    });
});
