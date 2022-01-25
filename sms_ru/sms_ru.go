package sms_ru

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SMS interface {
	Phone() string
	Message() string
	ID() string
}
type APIHandler struct {
	apiID   string
	baseUrl string
	Client  *http.Client
	Test    bool
}
type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMSID      string `json:"sms_id"`
	Cost       string `json:"cost"`
	StatusText string `json:"status_text"`
}
type ApiResponse struct {
	Status     string              `json:"status"`
	StatusCode int                 `json:"status_code"`
	SMS        map[string]Response `json:"sms"`
	Balance    float64             `json:"balance"`
}

type APIError struct {
	Err       string
	ErrCode   string
	ReqStatus string
}

func (err APIError) Error() string {
	return err.Err
}

func NewAPIHandler(apiID string, client *http.Client, test bool) *APIHandler {
	return &APIHandler{apiID: apiID, baseUrl: "https://sms.ru", Client: client, Test: test}
}
func (api *APIHandler) setBaseURL(baseUrl string) {
	api.baseUrl = baseUrl
}

func (api APIHandler) SendSMSList(smsList *[]SMS, json bool) (string, error) {
	var err error
	err = APIError{}
	clearRes := ""
	if len(*smsList) == 0 {
		return clearRes, errors.New("sms list is empty")
	}
	smsFromData := url.Values{}
	for _, msg := range *smsList {
		key := "to[" + msg.Phone() + "]"
		smsFromData.Add(key, msg.Message())
	}
	url := api.baseUrl + "/sms/send"
	req, err := http.NewRequest("POST", url, strings.NewReader(smsFromData.Encode()))
	if err != nil {
		return clearRes, err
	}
	//req.Header.Set("content-type", "application/x-www-form-urlencoded")
	q := req.URL.Query()
	q.Add("api_id", api.apiID)
	if json {
		q.Add("json", "1")
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = q.Encode()
	resp, err := api.Client.Do(req)
	if err != nil {
		return clearRes, err
	}
	aoiAnswer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clearRes, err
	}
	return string(aoiAnswer), nil
}
