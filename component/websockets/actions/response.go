package actions

import "encoding/json"

type JSONResponse struct {
	ErrorMessage string
	Response     interface{}
	ResponseData []byte
}

func StandardResponse(x interface{}) *JSONResponse {
	data, err := json.Marshal(x)
	if err != nil {
		return nil
	}
	return &JSONResponse{
		ErrorMessage: "",
		Response:     x,
		ResponseData: data,
	}
}
