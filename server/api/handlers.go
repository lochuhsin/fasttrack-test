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
		return
	}
	payload := internal.CreateUser{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&payload)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(RequestError{Msg: "Invalid payload"})
		resp.Write(b)
		return
	}

	if payload.Name == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "Invalid payload, missing name field"}.Bytes())
		return
	}
	// Not sure if we need to abstract this layer, maybe refactor it ?
	ok := internal.GetUserGroup().Create(*payload.Name)
	if !ok {
		resp.WriteHeader(http.StatusConflict)
		resp.Write(RequestError{Msg: "Invalid payload, missing name field"}.Bytes())
		return
	}
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("ok!"))
}

func ListAllQuestions(resp http.ResponseWriter, req *http.Request) { // pagination ?

}

func CreateRecord(resp http.ResponseWriter, req *http.Request) {
	reader := req.Body
	if reader == nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	payload := internal.CreateRecord{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(RequestError{Msg: "Invalid payload"})
		resp.Write(b)
		return
	}
	if payload.Name == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "Invalid payload, missing name field"}.Bytes())
		return
	}
	if payload.Answers == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "Invalid payload, missing answers field"}.Bytes())
		return
	}

	/**
	 * This should be guarded by mutex or something, or maybe channel ?
	 */
	ok := internal.GetUserGroup().Exists(*payload.Name)
	if !ok {
		resp.WriteHeader(http.StatusConflict)
		resp.Write(RequestError{Msg: "User not exists, something went wrong ?"}.Bytes())
		return
	}

	/** TODO:
	 * Calculate the score
	 *
	 * and
	 *
	 * think the design of concurrency
	 */

	// Maybe provide something like create or update ?
	ok = internal.GetSubmitRecord().Create(*payload.Name, 10)
	if !ok {
		resp.WriteHeader(http.StatusConflict)
		resp.Write(RequestError{Msg: "Record from this user was created"}.Bytes())
		return
	}
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("ok!"))
}
