package service

import (
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/simple-datasource-backend/pkg/mqtt"
	"net/http"
)

func (s *dataSourceInstance) Settings() (mqtt.MqttConfigurations, error) {
	return newDataSourceSettings(s.settings)
}

func newDataSourceSettings(instanceSettings backend.DataSourceInstanceSettings) (mqtt.MqttConfigurations, error) {
	var settings mqtt.MqttConfigurations
	if err := json.Unmarshal(instanceSettings.JSONData, &settings); err != nil {
		return mqtt.MqttConfigurations{}, err
	}
	if val, ok := instanceSettings.DecryptedSecureJSONData["password"]; ok {
		settings.Password = val
	}
	return settings, nil
}

func (ds *MqttDatasource) GetSettingsFromCtx(r *http.Request) (*backend.PluginContext, *mqtt.MqttConfigurations, error){
	pluginCtx := httpadapter.PluginConfigFromContext(r.Context())
	instance, err := ds.getInstance(pluginCtx)
	settings, err := instance.Settings()
	if err != nil {
		return nil, nil, err
	}
	return &pluginCtx, &settings, nil
}
