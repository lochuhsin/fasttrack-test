package api

import (
	"bytes"
	"encoding/json"
	"fasttrack-server/internal"
	"fmt"
	"net/http"
	"strconv"
)

func createUser(resp http.ResponseWriter, req *http.Request) {
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
		resp.Write(RequestError{Msg: "Invalid payload"}.Bytes())
		return
	}

	if payload.Name == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "Invalid payload, missing name field"}.Bytes())
		return
	}
	ok := internal.GetUserGroup().Create(*payload.Name)
	if !ok {
		resp.WriteHeader(http.StatusConflict)
		resp.Write(RequestError{Msg: "Invalid payload, missing name field"}.Bytes())
		return
	}
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("ok!"))
}

func listQuestions(resp http.ResponseWriter, req *http.Request) { // pagination ?
	query := req.URL.Query()
	limit := 10
	offset := 0
	if len(query["limit"]) != 0 {
		l, err := strconv.Atoi(query["limit"][0])
		if err != nil || l < 0 {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write(RequestError{Msg: "Invalid query limit parameters"}.Bytes())
			return
		}
		limit = l
	}

	if len(query["offset"]) != 0 {
		o, err := strconv.Atoi(query["offset"][0])
		if err != nil || o < 0 {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write(RequestError{Msg: "Invalid query offset parameters"}.Bytes())
			return
		}
		offset = o
	}

	qs, status := internal.GetQuestionDatabase().List(limit, offset)
	if !status {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "Invalid query parameters"}.Bytes())
		return
	}

	questions := make([]internal.Question, 0, len(qs))
	for _, q := range qs {
		questions = append(questions, internal.Question{
			Q:       q.Q,
			Options: q.Options,
		})
	}
	next := limit + offset + 1
	if next >= internal.GetQuestionDatabase().Count() {
		next = -1
	}

	var response *internal.ListQuestions
	if next < internal.GetQuestionDatabase().Count() {
		response = &internal.ListQuestions{
			Questions: questions,
			Next:      &next,
		}
	} else {
		response = &internal.ListQuestions{
			Questions: questions,
		}
	}

	b, err := json.Marshal(response)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write(RequestError{Msg: "Question Database error"}.Bytes())
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write(b)
}

func createRecord(resp http.ResponseWriter, req *http.Request) {
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
		resp.Write(RequestError{Msg: "Invalid payload"}.Bytes())
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
	count, err := internal.CreateRecordWithName(*payload.Name, payload.Answers)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: err.Error()}.Bytes())
		return
	}

	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte(fmt.Sprintf("You scored %v out of %v problems correctly", count, internal.GetQuestionDatabase().Count())))
}

func getPercentile(resp http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	if len(query["name"]) == 0 {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "missing name query"}.Bytes())
		return
	}

	val, ok := internal.GetRecords().GetPercentile(query["name"][0])
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(RequestError{Msg: "wrong name"}.Bytes())
		return
	}
	var builder bytes.Buffer
	builder.WriteString("You were better than ")
	builder.WriteString(strconv.Itoa(val))
	builder.WriteString("% of all quizzers")
	resp.WriteHeader(http.StatusOK)
	resp.Write(builder.Bytes())
}
