package sms_service

import (
	"context"
	"github.com/arkoil/sms_service/internal/background"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"
)

type Application struct {
	ErrorLog        *log.Logger
	InfoLog         *log.Logger
	CTX             *context.Context
	DB              *redis.Client
	APIsms          background.SMSAPIHandler
	KeyPrefixToSend string
	KeyPrefixBeSend string
	SendInterval    int
	StorePeriod     time.Duration
	JobsWG          *sync.WaitGroup
}

func (a *Application) RunJobs() {
	sendJob := background.NewSMSSendJob(
		a.DB,
		a.CTX,
		a.APIsms,
		a.InfoLog,
		a.ErrorLog,
		background.SetKeyPrefixToSend(a.KeyPrefixToSend),
		background.SetKeyPrefixBeSend(a.KeyPrefixBeSend),
		background.SetInterval(a.SendInterval),
		background.SetStorePeriod(a.StorePeriod),
	)
	go sendJob.Task(a.JobsWG)
}
