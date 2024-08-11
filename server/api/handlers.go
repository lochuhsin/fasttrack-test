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

func listQuestions(resp http.ResponseWriter, req *http.Request) { // pagination ?
	/**
	 * add pagination
	 */
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

	// Not sure if we need to abstract this layer, maybe refactor it ?
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
	questions, _ := internal.GetQuestionDatabase().List(-1, -1)
	count := 0
	for _, answer := range payload.Answers {
		qIndex, ans := answer.QuestionId, answer.Answer
		if questions[*qIndex].A == *ans {
			count++
		}

	}

	// Maybe provide something like create or update ?
	// Not sure if we need to abstract this layer, maybe refactor it ?
	ok = internal.GetRecords().Create(*payload.Name, count)
	if !ok {
		resp.WriteHeader(http.StatusConflict)
		resp.Write(RequestError{Msg: "Record from this user was created"}.Bytes())
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
	builder.WriteString(strconv.Itoa(val * 100))
	builder.WriteString("% of all quizzers")
	resp.WriteHeader(http.StatusOK)
	resp.Write(builder.Bytes())
}
