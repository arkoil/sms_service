package background

import (
	"github.com/arkoil/sms_service/internal/store"
	"log"
	"sync"
	"time"
)

type SMSSendJob struct {
	DB         *store.DB
	ApiHandler SMSAPIHandler
	interval   time.Duration
	infLog     *log.Logger
	errLog     *log.Logger
}

type SMSJobOption func(j SMSSendJob) SMSSendJob

// interfaces for use api

type SMSAPIResponseItem interface {
	GetID() string
	GetStatus() string
}
type SMSAPIResponse interface {
	GetStatus() string
	RowResponse() string
	Items() []SMSAPIResponseItem
}
type SMSAPIHandler interface {
	SendSMSList([]store.SMS) (SMSAPIResponse, error)
	SendSMS(store.SMS) (SMSAPIResponse, error)
}

// Set options functions

// SetInterval set option
func SetInterval(s int) SMSJobOption {
	return func(j SMSSendJob) SMSSendJob {
		j.interval = time.Duration(s) * time.Second
		return j
	}
}

// NewSMSSendJob create object
func NewSMSSendJob(db *store.DB, api SMSAPIHandler, infLog *log.Logger, errLog *log.Logger, opt ...SMSJobOption) *SMSSendJob {
	job := SMSSendJob{
		DB:         db,
		ApiHandler: api,
		interval:   60 * time.Second,
		infLog:     infLog,
		errLog:     errLog,
	}
	for _, o := range opt {
		job = o(job)
	}
	return &job
}

// Task function for background send smd
func (j *SMSSendJob) Task(wg *sync.WaitGroup) {
	defer j.Task(wg)
	time.Sleep(j.interval)
	j.infLog.Print("SMS sending")
	wg.Add(1)
	defer wg.Done()
	smsList, err := j.DB.GetSMSList()
	if err != nil {
		j.errLog.Fatal(err)
	}
	result, err := j.ApiHandler.SendSMSList(smsList)
	if err != nil {
		j.errLog.Fatal(err)
	}
	j.infLog.Printf("Response status %s", result.GetStatus())
	for _, item := range result.Items() {
		err = j.DB.SetSMSSent(item.GetID(), item.GetStatus())
		if err != nil {
			j.errLog.Fatal(err)
		}
	}
	j.infLog.Print("End SMS sending")
}
