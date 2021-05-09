package mqtt

type PublishRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type SubscribeRequest struct {
	Topic string `json:"topic"`
}

type UnsubscribeRequest struct {
	Topic string `json:"topic"`
}
