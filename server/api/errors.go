package api

import (
	"encoding/json"
	"net/http"
)

type RequestError struct {
	Msg string `json:"msg"`
}

func (e RequestError) Bytes() []byte {
	b, _ := json.Marshal(e)
	return b
}

func writeErrorRequest(resp http.ResponseWriter, msg string, code int) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write(RequestError{Msg: "Invalid payload"}.Bytes())
}
