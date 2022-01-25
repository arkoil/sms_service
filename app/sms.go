package app

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
	key := "sms:sending:" + string(requestID)
	p["key"] = key
	jsonP, _ := json.Marshal(p)
	err := a.SmsListDB.Set(*a.CTX, key, jsonP, 0).Err()
	if err != nil {
		a.ErrorLog.Print(err)
		p["status"] = "Error"
		p["error"] = err.Error()
	} else {
		a.InfoLog.Print("Смс добавлена в очередь")
		p["status"] = "OK"
	}
	json.NewEncoder(w).Encode(p)
}
