package mqtt

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	cmap "github.com/orcaman/concurrent-map"
	"time"
)

type MqttConfigurations struct {
	Broker   string `json:"broker"`
	Port     string `json:"port"`
	ClientId string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MqttClient struct {
	MqttConfig *MqttConfigurations
	Client     mqtt.Client
}

var queue cmap.ConcurrentMap

func NewMqttClient(config *MqttConfigurations) *MqttClient {
	queue = cmap.New()
	return &MqttClient{
		MqttConfig: config,
		Client:     makeMqttClient(config),
	}
}

func makeMqttClient(config *MqttConfigurations) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", config.Broker, config.Port))
	opts.SetClientID(config.ClientId)
	if config.Username != "" && config.Password != "" {
		opts.SetUsername(config.Username)
		opts.SetPassword(config.Password)
	}
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnReconnecting = reconnectHandler
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.ResumeSubs = true
	opts.KeepAlive = 10
	opts.AutoReconnect = false
	opts.CleanSession = true
	return mqtt.NewClient(opts)
}

func cb (exists bool, valueInMap interface{}, newValue interface{}) interface{} {
	nv := newValue.(Message)
	if !exists {
		return []Message{nv}
	}
	res := valueInMap.([]Message)
	return append(res, nv)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())
	log.DefaultLogger.Info("Message Handler", "Payload", payload, "Topic", topic)
	data := Message{
		Payload: payload,
		Topic: topic,
		Timestamp: time.Now().UTC(),
	}
	log.DefaultLogger.Info("Message Handler", "Data", data)
	queue.Upsert(topic, data, cb)
}

var reconnectHandler mqtt.ReconnectHandler = func(client mqtt.Client, options *mqtt.ClientOptions) {
	log.DefaultLogger.Info("Reconnect Handler", "Connected", fmt.Sprintf("%v", client.IsConnected()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.DefaultLogger.Info("Connect Handler", "Connected", fmt.Sprintf("%v", client.IsConnected()))
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.DefaultLogger.Info("Connection Lost Handler", "Error", err.Error())
}

func (mc MqttClient) Connect() error {
	if token := mc.Client.Connect(); token.Wait() && token.Error() != nil {
		return errors.New("Error during connect: " + token.Error().Error())
	}
	return nil
}

func (mc MqttClient) Subscribe(topic string) error {
	token := mc.Client.Subscribe(topic, 1, nil)
	token.Wait()
	if token.Error() != nil {
		return errors.New("Error during subscribe: " + token.Error().Error())
	}
	return nil
}

func (mc MqttClient) Unsubscribe(topic string) error {
	token := mc.Client.Unsubscribe(topic)
	token.Wait()
	if token.Error() != nil {
		return errors.New("Error during unsubscribe: " + token.Error().Error())
	}
	return nil
}

func (mc MqttClient) Publish(topic string, message string) error {
	token := mc.Client.Publish(topic, 1, false, message)
	token.Wait()
	if token.Error() != nil {
		return errors.New("Error during publish: " + token.Error().Error())
	}
	return nil
}

func (mc MqttClient) Disconnect() {
	mc.Client.Disconnect(250)
}

func (mc MqttClient) IsConnected() bool {
	return mc.Client.IsConnected()
}

func (mc MqttClient) GetData(topic string) []Message {
	var result []Message
	if queue != nil {
		data, ok := queue.Get(topic)
		listItems, check := data.([]Message)
		if ok && check {
			log.DefaultLogger.Info("GetData", "Payload", listItems)
			result = listItems
		}
	}
	return result
}

func (mc MqttClient) DeleteData(topic string) {
	queue.Remove(topic)
}

// IsSameConnection This allow to check if another client is connected to the broker with the same client identifier
// because this is a cause of random disconnections
func (mc MqttClient) IsSameConnection(settings *MqttConfigurations) bool {
	if settings.Broker == mc.MqttConfig.Broker && settings.Port == mc.MqttConfig.Port && settings.ClientId == mc.MqttConfig.ClientId {
		return true
	}
	return false
}