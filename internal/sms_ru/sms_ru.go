package sms_ru

import (
	"encoding/json"
	"github.com/arkoil/sms_service/internal/background"
	"github.com/arkoil/sms_service/internal/store"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type APIHandler struct {
	Client       *http.Client
	errLog       *log.Logger
	infLog       *log.Logger
	apiID        string
	baseUrl      string
	test         bool
	jsonResponse bool
}

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMSID      string `json:"sms_id"`
	Cost       string `json:"cost"`
	StatusText string `json:"status_text"`
}
type ApiResponse struct {
	Status      string              `json:"status"`
	StatusCode  int                 `json:"status_code"`
	SMSItems    map[string]Response `json:"sms"`
	Balance     float64             `json:"balance"`
	rowResponse string
}

type APIError struct {
	Err       string
	ErrCode   string
	ReqStatus string
}

func (a ApiResponse) GetStatus() string {
	return a.Status
}
func (a ApiResponse) RowResponse() string {
	return a.rowResponse
}
func (a ApiResponse) Items() []background.SMSAPIResponseItem {
	items := make([]background.SMSAPIResponseItem, 0, len(a.SMSItems))
	for _, item := range a.SMSItems {
		items = append(items, item)
	}
	return items
}

func (r Response) GetStatus() string {
	return r.Status
}
func (r Response) GetID() string {
	return r.SMSID
}

func (err APIError) Error() string {
	return err.Err
}

func NewAPIHandler(apiID string, client *http.Client, infLog *log.Logger, errLog *log.Logger, opt ...Options) APIHandler {
	api := APIHandler{
		apiID:   apiID,
		baseUrl: "https://sms.ru",
		Client:  client,
	}
	for _, o := range opt {
		api = o(api)
	}
	return api
}
func (api APIHandler) SendSMS(sms store.SMS) (background.SMSAPIResponse, error) {
	list := make([]store.SMS, 0, 1)
	list = append(list, sms)
	return api.SendSMSList(list)
}
func (api APIHandler) SendSMSList(smsList []store.SMS) (background.SMSAPIResponse, error) {
	var err error
	err = APIError{}
	response := ApiResponse{}
	if len(smsList) == 0 {
		api.infLog.Println("sms list is empty")
		return response, nil
	}
	smsFromData := url.Values{}
	for _, msg := range smsList {
		key := "to[" + msg.Phone() + "]"
		smsFromData.Add(key, msg.Message())
	}
	url := api.baseUrl + "/sms/send"
	req, err := http.NewRequest("POST", url, strings.NewReader(smsFromData.Encode()))
	if err != nil {
		return response, err
	}
	q := req.URL.Query()
	q.Add("api_id", api.apiID)
	if api.jsonResponse {
		q.Add("json", "1")
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = q.Encode()
	resp, err := api.Client.Do(req)
	if err != nil {
		return response, err
	}
	apiAnswer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	response.rowResponse = string(apiAnswer)
	err = json.Unmarshal(apiAnswer, &response)
	if err != nil {
		return response, err
	}
	return response, err
}
