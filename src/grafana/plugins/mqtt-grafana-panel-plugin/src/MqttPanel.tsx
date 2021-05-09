import React from 'react';
import { PanelProps } from '@grafana/data';
import { useState } from 'react';
import Response from './Response';
import { Button, Form } from './Components';
import { GetDataResponse, SimpleOptions } from 'types';
import { DataSourceWithBackend } from '@grafana/runtime';
import { DataSourceJsonData, DataQuery, DataSourceInstanceSettings } from '@grafana/data';

import * as c from './constants';
import './style.css';

interface Props extends PanelProps<SimpleOptions> {}

class MyDataSource extends DataSourceWithBackend<DataQuery, DataSourceJsonData> {
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceJsonData>) {
    super(instanceSettings);
  }
  // custom methods
}

export const MqttPanel: React.FC<Props> = ({ options, data, width, height, replaceVariables }) => {
  const { datasource, topic, color_button, color_text, operation } = options;
  const [payload, setPayload] = useState('');
  const [response, setResponse] = useState('');
  const [init, setInit] = useState(true);
  const [mqttData, setData] = useState('');

  const settings = {
    id: datasource,
  } as DataSourceInstanceSettings<DataSourceJsonData>;
  const ds = new MyDataSource(settings);
  const url = operation;

  /*if (init && operation === c.connectOp) {
    handleConnect();
    setInit(false);
  }*/

  if (init && operation === c.getDataOp) {
    handleGetData(topic);
    setInit(false);
  }

  function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  function handleDeleteData() {
    const req_url = url + '/' + topic + '/delete';
    console.log('Request: ' + req_url);
    ds.postResource(req_url).then((resp) => {
      setData('');
    });
  }

  function handleGetData(topic: any) {
    const req_url = url + '/' + topic;
    console.log('Request: ' + req_url);

    (async () => {
      while (true) {
        ds.getResource(req_url).then((resp) => {
          let list = resp.response;
          const listItems = list.map((elem: GetDataResponse) => <li key={elem.toString()}>{elem.payload}</li>);
          setData(listItems);
        });
        await delay(5000);
      }
    })();
  }

  function handleSubscribe() {
    const body = { topic: topic };
    console.log('Request: ' + url + ' Body: ' + JSON.stringify(body));

    ds.postResource(url, body).then((resp) => {
      if (resp.err !== '') {
        setResponse(resp.err);
        return;
      }
      setResponse(resp.response);
    });
  }

  function handleUnsubscribe() {
    const body = { topic: topic };
    console.log('Request: ' + url + ' Body: ' + JSON.stringify(body));

    ds.postResource(url, body).then((resp) => {
      if (resp.err !== '') {
        setResponse(resp.err);
        return;
      }
      setResponse(resp.response);
    });
  }

  /*function handleCheckConnection() {
    const url = c.isConnectedOp;
    console.log('Request: ' + url);

    ds.getResource(url).then((resp) => {
      if (resp.err !== '') {
        setResponse(resp.err);
        return;
      }
      setResponse(resp.response);
    });
  }*/

  function handleMessage(e: any) {
    setPayload(e.target.value);
  }

  function handlePublish() {
    const body = { topic: topic, message: payload };
    console.log('Request: ' + url + ' body: ' + JSON.stringify(body));

    ds.postResource(url, body).then((resp) => {
      if (resp.err !== '') {
        setResponse(resp.err);
        return;
      }
      setResponse(resp.response);
    });
  }

  function handleDisconnect() {
    const body = { topic: topic };
    console.log('Request: ' + url + ' Body: ' + JSON.stringify(body));

    ds.postResource(url).then((resp) => {
      if (resp.err !== '') {
        setResponse(resp.err);
        return;
      }
      setResponse(resp.response);
    });
  }

  function handleConnect() {
    console.log('Request: ' + url);

    ds.postResource(url).then((resp) => {
      if (resp.err !== '') {
        setResponse(resp.err);
        return;
      }
      setResponse(resp.response);
    });
  }
  switch (operation) {
    case c.connectOp:
      return (
        <div>
          <div className="centerFlex">
            <Button
              title={c.connectName}
              backgroundcolor={color_button}
              textcolor={color_text}
              classname={c.buttonClass}
              handle={handleConnect}
            ></Button>
          </div>
          {response !== null ? <Response value={response}></Response> : ''}
        </div>
      );
    case c.disconnectOp:
      return (
        <div>
          <div className="centerFlex">
            <Button
              title={c.disconnectName}
              backgroundcolor={color_button}
              textcolor={color_text}
              classname={c.buttonClass}
              handle={handleDisconnect}
            ></Button>
          </div>
          {response !== null ? <Response value={response}></Response> : ''}
        </div>
      );
    case c.subOp:
      return (
        <div>
          <div className="centerFlex">
            <Button
              title={c.subName}
              backgroundcolor={color_button}
              textcolor={color_text}
              classname={c.buttonClass}
              handle={handleSubscribe}
            ></Button>
          </div>
          {response !== null ? <Response value={response}></Response> : ''}
        </div>
      );
    case c.unsubscribeOp:
      return (
        <div>
          <div className="centerFlex">
            <Button
              title={c.unsubscribeName}
              backgroundcolor={color_button}
              textcolor={color_text}
              classname={c.buttonClass}
              handle={handleUnsubscribe}
            ></Button>
          </div>
          {response !== null ? <Response value={response}></Response> : ''}
        </div>
      );
    case c.publishOp:
      return (
        <div>
          <div className="centerFlex">
            <Form type="text" value={payload} handle={handleMessage}></Form>
            <Button
              title={c.publishName}
              backgroundcolor={color_button}
              textcolor={color_text}
              classname={c.buttonClass}
              handle={handlePublish}
            ></Button>
          </div>
          {response !== null ? <Response value={response}></Response> : ''}
        </div>
      );
    case c.getDataOp:
      return (
        <div>
          <Response value={c.getDataTitle + topic}></Response>
          <Response value={mqttData}></Response>
          <Button
            title={c.deleteName}
            backgroundcolor={color_button}
            textcolor={color_text}
            classname={c.buttonClass}
            handle={handleDeleteData}
          ></Button>
        </div>
      );
    default:
      return (
        <div>
          <h1>Unknown</h1>
        </div>
      );
  }
};
