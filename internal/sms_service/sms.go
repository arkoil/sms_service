package sms_service

import (
	"encoding/json"
	"net/http"
)

func (a Application) errorResponse(errorMsg string, w http.ResponseWriter) {
	a.ErrorLog.Print(errorMsg)
	response := make(map[string]string)
	response["status"] = "Error"
	response["error"] = errorMsg
	a.jsonResponse(response, w)
}
func (a Application) jsonResponse(response map[string]string, w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		a.ErrorLog.Fatal(err)
	}

}
func (a Application) SendSMS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	message := make(map[string]string)
	response := make(map[string]string)
	message["number"] = r.URL.Query().Get("number")
	message["message"] = r.URL.Query().Get("message")
	message["request_id"] = r.URL.Query().Get("request_id")
	if message["number"] == "" {
		a.errorResponse("number is empty", w)
		return
	}
	if message["message"] == "" {
		a.errorResponse("message is empty", w)
		return

	}
	if message["request_id"] == "" {
		a.errorResponse("request_id is empty", w)
		return
	}
	key := a.KeyPrefixToSend + message["request_id"]
	jsonMessage, _ := json.Marshal(message)
	err := a.DB.Set(*a.CTX, key, jsonMessage, 0).Err()
	if err != nil {
		a.errorResponse(err.Error(), w)
	}
	a.InfoLog.Print("SMS added to the queue")
	response["status"] = "OK"
	response["result"] = "SMS added to the queue"
	a.jsonResponse(response, w)

}

func (a Application) CheckSMS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	requestID := r.URL.Query().Get("request_id")
	if requestID == "" {
		a.errorResponse("request_id is empty", w)
		return
	}
	key := a.KeyPrefixBeSend + requestID
	redisResult, err := a.DB.Get(*a.CTX, key).Result()
	if err != nil {
		a.errorResponse(err.Error(), w)
		return
	}
	a.InfoLog.Print("SMS added in turn")
	response["status"] = "OK"
	response["result"] = redisResult
	a.jsonResponse(response, w)
}
