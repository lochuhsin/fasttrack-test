package main

import (
	"fasttrack-test/server/api"
	"fasttrack-test/server/internal"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Init() {
	internal.InitQuestionDatabase()
	internal.InitUserGroup()
	internal.InitSubmitRecord()
}

func main() {
	/**
	 * move port to environment
	 */
	Init()
	mux := http.NewServeMux()
	api.Register(mux)
	log.Info("Running server ...")
	http.ListenAndServe(":8000", mux)

}
