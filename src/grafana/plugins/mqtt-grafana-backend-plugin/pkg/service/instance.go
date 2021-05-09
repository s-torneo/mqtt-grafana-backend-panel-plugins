package service

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

type dataSourceInstance struct {
	settings   backend.DataSourceInstanceSettings
}

func newDataSourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &dataSourceInstance{
		settings: settings,
	}, nil
}

func (ds *MqttDatasource) getInstance(ctx backend.PluginContext) (*dataSourceInstance, error) {
	instance, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return instance.(*dataSourceInstance), nil
}

func (s *dataSourceInstance) Dispose() {
	// Called before creating a a new instance to allow plugin authors
	// to cleanup.window.onload = function () {window.location.reload()}
}