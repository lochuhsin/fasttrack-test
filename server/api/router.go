package api

import "net/http"

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("POST /records", createRecord)
	mux.HandleFunc("GET /questions", listQuestions)
	mux.HandleFunc("GET /percentile", getPercentile)
}
