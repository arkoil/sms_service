package sms_service

import (
	"github.com/gorilla/mux"
)

func (a Application) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/sms/send", a.SendSMS).Methods("POST")
	router.HandleFunc("/sms/check", a.CheckSMS).Methods("GET")
	return router
}
