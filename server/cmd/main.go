package main

import (
	"fasttrack-server/api"
	"fasttrack-server/internal"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Init() {
	internal.InitQuestionDatabase()
	internal.InitUserGroup()
	internal.InitRecords()
}

func main() {
	/**
	 * move port to environment
	 */
	Init()
	mux := http.NewServeMux()
	api.Register(mux)
	log.Info("Running server ... on 8000")
	http.ListenAndServe(":8000", mux)

}
