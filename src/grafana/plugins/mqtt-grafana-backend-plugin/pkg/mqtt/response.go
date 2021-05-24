package mqtt

import "time"

type Response struct {
	Res string `json:"response"`
	Err string `json:"err"`
}

type Message struct {
	Topic     string    `json:"topic"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"ts"`
}

type GetDataResponse struct {
	Res []Message `json:"response"`
	Err string    `json:"err"`
}

func MakeResponse(response string, err error) Response {
	if err != nil {
		return Response{
			Res: "",
			Err: err.Error(),
		}
	}
	return Response{
		Res: response,
		Err: "",
	}
}
