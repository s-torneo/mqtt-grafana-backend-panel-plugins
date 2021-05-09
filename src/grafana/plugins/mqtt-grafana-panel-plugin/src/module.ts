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
      path: 'topic',
      name: 'Topic',
      description: 'MQTT Topic',
      defaultValue: 'test',
    })
    .addColorPicker({
      path: 'color_text',
      name: 'Text Color of button',
      description: '',
      defaultValue: '',
    })
    .addColorPicker({
      path: 'color_button',
      name: 'Button Color',
      description: 'Color of the button',
      defaultValue: '',
    })
    .addTextInput({
      path: 'datasource',
      name: 'Datasource Id',
      description: 'Id of the datasource',
      defaultValue: '1',
    });
});
