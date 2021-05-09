package mqtt

type Response struct {
	Res string `json:"response"`
	Err string `json:"err"`
}

type DataResponse struct {
	Payload string `json:"payload"`
}

type GetDataResponse struct {
	Res []DataResponse `json:"response"`
	Err string        `json:"err"`
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

func MakeGetDataResponse(data []string, err error) GetDataResponse {
	if err != nil {
		return GetDataResponse{
			Res: nil,
			Err: err.Error(),
		}
	}
	var response []DataResponse
	for _, d := range data {
		response = append(response, DataResponse{Payload: d})
	}
	return GetDataResponse{
		Res: response,
		Err: "",
	}
}
