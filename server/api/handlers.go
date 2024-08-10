package api

import (
	"encoding/json"
	"fasttrack-test/server/internal"
	"net/http"
)

func CreateUser(resp http.ResponseWriter, req *http.Request) {
	reader := req.Body
	if reader == nil {
		resp.WriteHeader(http.StatusBadRequest)
	}
	payload := internal.CreateUser{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&payload)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(RequestError{Msg: "Invalid payload"})
		resp.Write(b)
	}

	if payload.Name == nil {
		resp.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(RequestError{Msg: "Invalid payload, missing name field"})
		resp.Write(b)
	}
	// Not sure if we need to abstract this layer, maybe refactor it ?
}

func ListAllQuestions(resp http.ResponseWriter, req *http.Request) { // pagination ?

}

func CreateRecord(resp http.ResponseWriter, req *http.Request) {

}
