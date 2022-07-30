package sms_service

import (
	"github.com/gorilla/mux"
)

func (a Application) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/sms/send", a.SendSMS)
	router.HandleFunc("/sms/check", a.SendSMS)
	return router
}
