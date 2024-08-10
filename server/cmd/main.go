package main

import (
	"fasttrack-test/server/api"
	"fasttrack-test/server/internal"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Init() {
	internal.InitUserGroup()
	internal.InitSubmitRecord()
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", api.CreateUser)
	mux.HandleFunc("POST /records", api.CreateRecord)
	mux.HandleFunc("GET /questions", api.ListAllQuestions)
}

func main() {
	/**
	 * move port to environment
	 */
	Init()
	mux := http.NewServeMux()
	Register(mux)
	log.Info("Running server ...")
	http.ListenAndServe(":8000", mux)

}
