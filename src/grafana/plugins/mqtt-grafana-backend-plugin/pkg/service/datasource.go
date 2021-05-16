package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/simple-datasource-backend/pkg/mqtt"
	"net/http"
)

type MqttDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im         instancemgmt.InstanceManager
	mqttClient *mqtt.MqttClient
}

// NewDatasource returns datasource.ServeOpts.
func NewDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	mqttClient := mqtt.NewMqttClient(&mqtt.MqttConfigurations{})

	ds := &MqttDatasource{
		im:         im,
		mqttClient: mqttClient,
	}

	return datasource.ServeOpts{
		CheckHealthHandler:  ds,
		CallResourceHandler: httpadapter.New(ds),
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (ds *MqttDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	instance, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	settings, err := instance.Settings()
	if err != nil {
		log.DefaultLogger.Error("CheckHealth", err.Error())
		return nil, err
	}

	log.DefaultLogger.Info("CheckHealth", "PluginContext", req.PluginContext, "Settings", settings)

	if ds.mqttClient.IsSameConnection(&settings) {
		ds.mqttClient.Disconnect()
	}
	ds.mqttClient = mqtt.NewMqttClient(&settings)
	err = ds.mqttClient.Connect()
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Mqtt client is working",
	}, nil
}

// ServeHTTP is the main HTTP handler for serving resource calls
// https://github.com/grafana/github-datasource
func (ds *MqttDatasource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/connect").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.DefaultLogger.Info("ServeHTTP", "Request", "Connect")
			pluginCtx, settings, err := ds.GetSettingsFromCtx(r)
			if err != nil {
				log.DefaultLogger.Error("GetSettingsFromCtx", "Error", err.Error())
				//w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(mqtt.MakeResponse("", err))
				return
			}
			log.DefaultLogger.Info("PluginContext", pluginCtx, "PluginSettings", settings)
			if ds.mqttClient.IsSameConnection(settings) {
				ds.mqttClient.Disconnect()
			}
			ds.mqttClient = mqtt.NewMqttClient(settings)
			err = ds.mqttClient.Connect()
			if err != nil {
				log.DefaultLogger.Error(err.Error())
				//w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(mqtt.MakeResponse("", err))
				return
			}
			json.NewEncoder(w).Encode(mqtt.MakeResponse("Connected", err))
		})

	router.Methods(http.MethodPost).Path("/subscribe").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req mqtt.SubscribeRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			log.DefaultLogger.Info("ServeHTTP", "Request", "Subscribe", "Topic", req.Topic)
			err = ds.mqttClient.Subscribe(req.Topic)
			if err != nil {
				log.DefaultLogger.Error(err.Error())
				//w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(mqtt.MakeResponse("", err))
				return
			}
			json.NewEncoder(w).Encode(mqtt.MakeResponse(fmt.Sprintf("Subscribed to topic: %s", req.Topic), nil))
		})

	router.Methods(http.MethodPost).Path("/unsubscribe").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req mqtt.UnsubscribeRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			log.DefaultLogger.Info("ServeHTTP", "Request", "Unsubscribe", "Topic", req.Topic)
			err = ds.mqttClient.Unsubscribe(req.Topic)
			if err != nil {
				log.DefaultLogger.Error(err.Error())
				//w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(mqtt.MakeResponse("", err))
				return
			}
			json.NewEncoder(w).Encode(mqtt.MakeResponse(fmt.Sprintf("Unsubscribed from topic: %s", req.Topic), nil))
		})

	router.Methods(http.MethodPost).Path("/publish").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req mqtt.PublishRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			log.DefaultLogger.Info("ServeHTTP", "Request", "Publish", "Topic", req.Topic, "Message", req.Message)
			err = ds.mqttClient.Publish(req.Topic, req.Message)
			if err != nil {
				log.DefaultLogger.Error(err.Error())
				//w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(mqtt.MakeResponse("", err))
				return
			}
			json.NewEncoder(w).Encode(mqtt.MakeResponse(fmt.Sprintf("Published message: %s to topic: %s", req.Message, req.Topic), nil))
		})

	router.Methods(http.MethodPost).Path("/disconnect").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.DefaultLogger.Info("ServeHTTP", "Request", "Disconnect")
			ds.mqttClient.Disconnect()
			json.NewEncoder(w).Encode(mqtt.MakeResponse(fmt.Sprintf("Disconnected"), nil))
		})

	router.Methods(http.MethodGet).Path("/connection").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.DefaultLogger.Info("ServeHTTP", "Request", "ConnectionStatus")
			check := ds.mqttClient.IsConnected()
			json.NewEncoder(w).Encode(mqtt.MakeResponse(fmt.Sprintf("%v", check), nil))
		})

	router.Methods(http.MethodGet).Path("/data/{topic}").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			topic := vars["topic"]
			log.DefaultLogger.Info("ServeHTTP", "Request", "GetData", "Topic", topic)
			data := ds.mqttClient.GetData(topic)
			json.NewEncoder(w).Encode(mqtt.MakeGetDataResponse(data, nil))
		})

	router.Methods(http.MethodPost).Path("/data/{topic}/delete").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			topic := vars["topic"]
			log.DefaultLogger.Info("ServeHTTP", "Request", "DeleteData", "Topic", topic)
			ds.mqttClient.DeleteData(topic)
		})

	router.ServeHTTP(w, r)
}