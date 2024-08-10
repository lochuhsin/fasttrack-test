package api

import "encoding/json"

type RequestError struct {
	Msg string `json:"msg"`
}

func (e RequestError) Bytes() []byte {
	b, _ := json.Marshal(e)
	return b
}
