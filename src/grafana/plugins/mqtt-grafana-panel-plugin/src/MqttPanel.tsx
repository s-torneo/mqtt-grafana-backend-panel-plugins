import React from 'react';
import { PanelProps } from '@grafana/data';
import { useState, useRef } from 'react';
import Response from './Response';
import { Button, Form } from './Components';
import { GetDataResponse, SimpleOptions } from 'types';
import { DataSourceWithBackend } from '@grafana/runtime';
import { DataSourceJsonData, DataQuery, DataSourceInstanceSettings } from '@grafana/data';
import { CronJob } from 'cron';
import { Badge } from 'antd';

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
  const { datasource, mqttTopic, color_button, color_text, operation, publishMsg, buttonName } = options;
  const [message, setMessage] = useState('');
  const [response, setResponse] = useState('');
  const [init, setInit] = useState(true);
  const [mqttData, setData] = useState('');

  const settings = {
    id: datasource,
  } as DataSourceInstanceSettings<DataSourceJsonData>;
  const ds = new MyDataSource(settings);
  const url = operation;

  const lastTopic = useRef();
  let topic = replaceVariables(mqttTopic);
  // @ts-ignore
  lastTopic.current = topic;

  if (init && operation === c.getDataOp) {
    scheduleJob(handleGetData);
  }

  if (init && operation === c.isConnectedOp) {
    scheduleJob(handleCheckConnection);
  }

  function scheduleJob(func: any) {
    let cronJob = new CronJob('*/5 * * * * *', async () => {
      await func();
    });

    // Start job
    if (!cronJob.running) {
      cronJob.start();
    }
    setInit(false);
  }

  function handleDeleteData() {
    const req_url = url + '/' + topic + '/delete';
    console.log('Request: ' + req_url);
    ds.postResource(req_url).then((resp) => {
      setData('');
    });
  }

  async function handleGetData() {
    const req_url = url + '/' + lastTopic.current;
    console.log('Request: ' + req_url);
    ds.getResource(req_url).then((resp) => {
      let list = resp.response;
      if (list !== null) {
        const listItems = list.map((elem: GetDataResponse) => <li key={elem.toString()}>{elem.payload}</li>);
        setData(listItems);
      }
    });
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

  function handleCheckConnection() {
    const url = c.isConnectedOp;
    console.log('Request: ' + url);

    ds.getResource(url).then((resp) => {
      setResponse(resp.response);
    });
  }

  function handleMessage(e: any) {
    setMessage(e.target.value);
  }

  function handlePublish() {
    let payload = '';
    if (publishMsg !== undefined && publishMsg !== '') {
      payload = publishMsg;
    } else {
      payload = message;
    }
    const body = { topic: topic, message: payload };
    console.log('Request: ' + url + ' Body: ' + JSON.stringify(body));

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
              title={buttonName === undefined || buttonName === '' ? c.connectName : buttonName}
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
              title={buttonName === undefined || buttonName === '' ? c.disconnectName : buttonName}
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
              title={buttonName === undefined || buttonName === '' ? c.subName : buttonName}
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
              title={buttonName === undefined || buttonName === '' ? c.unsubscribeName : buttonName}
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
            {publishMsg === undefined || publishMsg === '' ? (
              <Form type="text" value={message} handle={handleMessage}></Form>
            ) : (
              ''
            )}
            <Button
              title={buttonName === undefined || buttonName === '' ? c.publishName : buttonName}
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
            title={buttonName === undefined || buttonName === '' ? c.deleteName : buttonName}
            backgroundcolor={color_button}
            textcolor={color_text}
            classname={c.buttonClass}
            handle={handleDeleteData}
          ></Button>
        </div>
      );
    case c.isConnectedOp:
      return (
        <div>
          <Response value="Connection Status"></Response>
          {response === 'true' ? (
            <Badge dot={true} status="success" text="Connected" />
          ) : (
            <Badge dot={true} status="error" text="Disconnected" />
          )}
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
