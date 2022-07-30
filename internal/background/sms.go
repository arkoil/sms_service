package background

import (
	"context"
	"encoding/json"
	"github.com/arkoil/sms_service/internal/sms"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"
)

type SMSSendJob struct {
	DB              *redis.Client
	CTX             *context.Context
	ApiHandler      SMSAPIHandler
	keyPrefixToSend string
	keyPrefixBeSend string
	storePeriod     time.Duration
	interval        time.Duration
	infLog          *log.Logger
	errLog          *log.Logger
}

type SMSJobOption func(j SMSSendJob) SMSSendJob

// interfaces for use api

type SMS interface {
	Phone() string
	Message() string
	ID() string
}
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
	SendSMSList([]SMS) (SMSAPIResponse, error)
	SendSMS(SMS) (SMSAPIResponse, error)
}

// Set options functions

// SetInterval set option
func SetInterval(s int) SMSJobOption {
	return func(j SMSSendJob) SMSSendJob {
		j.interval = time.Duration(s) * time.Second
		return j
	}
}

// SetStorePeriod set option
func SetStorePeriod(sp time.Duration) SMSJobOption {
	return func(j SMSSendJob) SMSSendJob {
		j.storePeriod = sp * time.Hour
		return j
	}
}

// SetKeyPrefixToSend set option
func SetKeyPrefixToSend(pref string) SMSJobOption {
	return func(j SMSSendJob) SMSSendJob {
		j.keyPrefixToSend = pref
		return j
	}
}

// SetKeyPrefixBeSend set option
func SetKeyPrefixBeSend(pref string) SMSJobOption {
	return func(j SMSSendJob) SMSSendJob {
		j.keyPrefixBeSend = pref
		return j
	}
}

// NewSMSSendJob create object
func NewSMSSendJob(db *redis.Client, ctx *context.Context, api SMSAPIHandler, infLog *log.Logger, errLog *log.Logger, opt ...SMSJobOption) *SMSSendJob {
	job := SMSSendJob{
		DB:              db,
		CTX:             ctx,
		ApiHandler:      api,
		keyPrefixToSend: "",
		keyPrefixBeSend: "",
		interval:        60 * time.Second,
		infLog:          infLog,
		errLog:          errLog,
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
	keys, err := j.DB.Keys(*j.CTX, j.keyPrefixToSend+"*").Result()
	if err != nil {
		j.errLog.Fatal(err)
	} else {

		smsList := make([]SMS, 0)
		for _, key := range keys {
			jsonItem, _ := j.DB.Get(*j.CTX, key).Result()
			item := sms.Item{}
			err := json.Unmarshal([]byte(jsonItem), &item)
			if err != nil {
				j.errLog.Print(err)
				continue
			}
			smsList = append(smsList, item)
			j.DB.Del(*j.CTX, key)

		}
		result, err := j.ApiHandler.SendSMSList(smsList)
		if err != nil {
			j.errLog.Fatal(err)
		}
		j.infLog.Printf("Response status %s", result.GetStatus())
		for _, item := range result.Items() {
			key := j.keyPrefixBeSend + item.GetID()
			err := j.DB.Set(*j.CTX, key, item.GetStatus(), j.storePeriod).Err()
			if err != nil {
				j.errLog.Fatal(err)
			}
		}
	}
	j.infLog.Print("End SMS sending")
}
