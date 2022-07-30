package sms_service

import (
	"encoding/json"
	"net/http"
)

func (a Application) SendSMS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := make(map[string]string)
	num := r.URL.Query().Get("number")
	msg := r.URL.Query().Get("message")
	requestID := r.URL.Query().Get("request_id")
	p["number"] = num
	p["message"] = msg
	p["request_id"] = requestID
	key := a.KeyPrefixToSend + requestID
	p["key"] = key
	jsonP, _ := json.Marshal(p)
	err := a.DB.Set(*a.CTX, key, jsonP, 0).Err()
	if err != nil {
		a.ErrorLog.Print(err)
		p["status"] = "Error"
		p["error"] = err.Error()
	} else {
		a.InfoLog.Print("SMS added to the queue")
		p["status"] = "OK"
		p["result"] = "SMS added to the queue"
	}
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		a.ErrorLog.Fatal(err)
	}
}

func (a Application) CheckSMS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := make(map[string]string)
	requestID := r.URL.Query().Get("request_id")
	p["request_id"] = requestID
	key := a.KeyPrefixBeSend + requestID
	redisResult, err := a.DB.Get(*a.CTX, key).Result()
	if err != nil {
		a.ErrorLog.Print(err)
		p["status"] = "Error"
		p["error"] = err.Error()
	} else {
		a.InfoLog.Print("SMS added in turn")
		p["status"] = "OK"
		p["result"] = redisResult
	}
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		a.ErrorLog.Fatal(err)
	}

}
