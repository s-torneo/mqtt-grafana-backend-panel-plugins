package mqtt

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	cmap "github.com/orcaman/concurrent-map"
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

/*func makeTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	//cert := "-----BEGIN CERTIFICATE-----\nMIIF3jCCA8agAwIBAgIQAf1tMPyjylGoG7xkDjUDLTANBgkqhkiG9w0BAQwFADCB\niDELMAkGA1UEBhMCVVMxEzARBgNVBAgTCk5ldyBKZXJzZXkxFDASBgNVBAcTC0pl\ncnNleSBDaXR5MR4wHAYDVQQKExVUaGUgVVNFUlRSVVNUIE5ldHdvcmsxLjAsBgNV\nBAMTJVVTRVJUcnVzdCBSU0EgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMTAw\nMjAxMDAwMDAwWhcNMzgwMTE4MjM1OTU5WjCBiDELMAkGA1UEBhMCVVMxEzARBgNV\nBAgTCk5ldyBKZXJzZXkxFDASBgNVBAcTC0plcnNleSBDaXR5MR4wHAYDVQQKExVU\naGUgVVNFUlRSVVNUIE5ldHdvcmsxLjAsBgNVBAMTJVVTRVJUcnVzdCBSU0EgQ2Vy\ndGlmaWNhdGlvbiBBdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIK\nAoICAQCAEmUXNg7D2wiz0KxXDXbtzSfTTK1Qg2HiqiBNCS1kCdzOiZ/MPans9s/B\n3PHTsdZ7NygRK0faOca8Ohm0X6a9fZ2jY0K2dvKpOyuR+OJv0OwWIJAJPuLodMkY\ntJHUYmTbf6MG8YgYapAiPLz+E/CHFHv25B+O1ORRxhFnRghRy4YUVD+8M/5+bJz/\nFp0YvVGONaanZshyZ9shZrHUm3gDwFA66Mzw3LyeTP6vBZY1H1dat//O+T23LLb2\nVN3I5xI6Ta5MirdcmrS3ID3KfyI0rn47aGYBROcBTkZTmzNg95S+UzeQc0PzMsNT\n79uq/nROacdrjGCT3sTHDN/hMq7MkztReJVni+49Vv4M0GkPGw/zJSZrM233bkf6\nc0Plfg6lZrEpfDKEY1WJxA3Bk1QwGROs0303p+tdOmw1XNtB1xLaqUkL39iAigmT\nYo61Zs8liM2EuLE/pDkP2QKe6xJMlXzzawWpXhaDzLhn4ugTncxbgtNMs+1b/97l\nc6wjOy0AvzVVdAlJ2ElYGn+SNuZRkg7zJn0cTRe8yexDJtC/QV9AqURE9JnnV4ee\nUB9XVKg+/XRjL7FQZQnmWEIuQxpMtPAlR1n6BB6T1CZGSlCBst6+eLf8ZxXhyVeE\nHg9j1uliutZfVS7qXMYoCAQlObgOK6nyTJccBz8NUvXt7y+CDwIDAQABo0IwQDAd\nBgNVHQ4EFgQUU3m/WqorSs9UgOHYm8Cd8rIDZsswDgYDVR0PAQH/BAQDAgEGMA8G\nA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQEMBQADggIBAFzUfA3P9wF9QZllDHPF\nUp/L+M+ZBn8b2kMVn54CVVeWFPFSPCeHlCjtHzoBN6J2/FNQwISbxmtOuowhT6KO\nVWKR82kV2LyI48SqC/3vqOlLVSoGIG1VeCkZ7l8wXEskEVX/JJpuXior7gtNn3/3\nATiUFJVDBwn7YKnuHKsSjKCaXqeYalltiz8I+8jRRa8YFWSQEg9zKC7F4iRO/Fjs\n8PRF/iKz6y+O0tlFYQXBl2+odnKPi4w2r78NBc5xjeambx9spnFixdjQg3IM8WcR\niQycE0xyNN+81XHfqnHd4blsjDwSXWXavVcStkNr/+XeTWYRUc+ZruwXtuhxkYze\nSf7dNXGiFSeUHM9h4ya7b6NnJSFd5t0dCy5oGzuCr+yDZ4XUmFF0sbmZgIn/f3gZ\nXHlKYC6SQK5MNyosycdiyA5d9zZbyuAlJQG03RoHnHcAP9Dc1ew91Pq7P8yF1m9/\nqS3fuQL39ZeatTXaw2ewh0qpKJ4jjv9cJ2vhsE/zB+4ALtRZh8tSQZXq9EfX7mRB\nVXyNWQKV3WKdwrnuWih0hKWbt5DHDAff9Yk2dDLWKMGwsAvgnEzDHNb842m1R0aB\nL6KCq9NjRHDEjf8tM7qtj3u1cIiuPhnPQCjY/MiQu12ZIvVS5ljFH4gxQ+6IHdfG\njjxDah2nGN59PRbxYvnKkKj9\n-----END CERTIFICATE-----\n"

	config := &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
	return config
}*/

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
	// Create tls.Config with desired tls properties
	//opts.TLSConfig = makeTLSConfig()
	return mqtt.NewClient(opts)
}

/*func makeMqttWSClient(){
	host := "test"
	opts := &mqtt.WebsocketOptions{
		Proxy: ,
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	conn, _ := mqtt.NewWebsocket(
		host,
		&tls.Config{InsecureSkipVerify: false},
		60 * time.Second,
		http.Header{},
		opts)
}*/

func cb (exists bool, valueInMap interface{}, newValue interface{}) interface{} {
	nv := newValue.(string)
	if !exists {
		return []string{nv}
	}
	res := valueInMap.([]string)
	return append(res, nv)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())
	log.DefaultLogger.Info("Message Handler", "Payload", payload, "Topic", topic)
	queue.Upsert(topic, payload, cb)
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

func (mc MqttClient) GetData(topic string) []string {
	var result []string
	if queue != nil {
		data, ok := queue.Get(topic)
		listItems, check := data.([]string)
		if ok && check {
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