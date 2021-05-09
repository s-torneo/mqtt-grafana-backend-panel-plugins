# mqtt-grafana-backend-panel-plugins

A backend plugin in Go and a panel plugin in React and Typescript. 
The backend makes a mqtt client that you connect to a broker according to the configurations set, receives http internal requests from panel plugin and gives the response to them in order that panel plugin can show the response.

You can configure one or more Mqtt client (Grafana backend datasource) and one or more Grafana panel plugin that send requests to the Mqtt client selected (at the moment you have to specify the datasource id in the panel plugin settings). 

In the panel plugin you can select one of the following operations:
- Connect
- Disconnect
- Subscribe
- Unsubscribe
- Publish
- GetData

in order to show a different interface according to the operation selected.
You can also insert the mqtt topic and choose the color of the button and of the text.

![backend-plugin](https://github.com/s-torneo/mqtt-grafana-backend-panel-plugins/blob/master/img/backend-plugin.png "backend-plugin")

![panel-plugin](https://github.com/s-torneo/mqtt-grafana-backend-panel-plugins/blob/master/img/panel-plugin.png "panel-plugin")

You can run it in the following way:

- DEMO_PORT={your_port} docker-compose up --build

This docker-compose start two containers, one with grafana and one with nginx.
An example of datasource (grafana backend plugin) and of dashboard is loaded directly by docker-compose.

You can access to the Grafana to the address: http://localhost:{your_port}.


To compile and build Grafana plugins:

- Go in the directory: /src/grafana/plugins/mqtt-grafana-backend-plugin and run:
  - mage -v (to build backend)
  - yarn install (to install frontend dependencies)
  - yarn build (to build frontend)
 
- Go in the directory: /src/grafana/plugins/mqtt-grafana-panel-plugin and run:
  - yarn install (to install frontend dependencies)
  - yarn build (to build frontend)
  
  
For the frontend developing you don't need to run DEMO_PORT={your_port} docker-compose up --build each time, but you only need to run the command descripted above to compile and build frontend, instead for backend code you need to run mage -v and run again DEMO_PORT={your_port} docker-compose up --build.

This is only a way to develop quickly Grafana plugins, if you want use the plugins, you need only to copy the directories /src/grafana/plugins/mqtt-grafana-backend-plugin and /src/grafana/plugins/mqtt-grafana-panel-plugin inside the plugins directory of your Grafana.
