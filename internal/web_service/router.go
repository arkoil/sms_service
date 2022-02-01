package web_service

import (
	"github.com/gorilla/mux"
)

func (a Application) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/sms/send", a.SendSMS)
	return router
}
