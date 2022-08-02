package sms_service

import (
	"context"
	"github.com/arkoil/sms_service/internal/background"
	"github.com/arkoil/sms_service/internal/store"
	"log"
	"sync"
)

type Application struct {
	ErrorLog     *log.Logger
	InfoLog      *log.Logger
	CTX          *context.Context
	DB           *store.DB
	APIsms       background.SMSAPIHandler
	SendInterval int
	JobsWG       *sync.WaitGroup
}

func (a *Application) RunJobs() {
	sendJob := background.NewSMSSendJob(
		a.DB,
		a.APIsms,
		a.InfoLog,
		a.ErrorLog,
		background.SetInterval(a.SendInterval),
	)
	go sendJob.Task(a.JobsWG)
}
