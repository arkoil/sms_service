package main

import (
	"flag"
	"github.com/arkoil/sms_service/internal/sms_ru"
	"github.com/arkoil/sms_service/internal/sms_service"
	"github.com/arkoil/sms_service/internal/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	port := flag.String("port", "7575", "set service port")
	redisURL := flag.String("rurl", "localhost:6379", "set redis url")
	redisPassword := flag.String("rpass", "", "set redis password")
	redisDB := flag.Int("rdb", 11, "set redis db")
	keyPrefixToSend := flag.String("prefto", "sms:tosend:", "sets the Redis key to send messages")
	keyPrefixBeSend := flag.String("prefbe", "sms:besend:", "sets the Redis key for sent messages")
	sendInterval := flag.Int("interval", 60, "interval between sending SMS in seconds")
	smsRUAPIID := flag.String("sms-ru-aoi-id", "", "set sms.ru apiID")
	storaPeriod := flag.Int("sp", 24*30, " set storage period of sent SMS history in hours")
	flag.Parse()

	//Variables
	serverPort := strings.Join([]string{":", *port}, "")
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	var wg sync.WaitGroup
	//Clean close
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	//Redis client
	rdb, err := store.NewDB(*redisURL, *redisPassword, *redisDB, *keyPrefixToSend, *keyPrefixBeSend, *storaPeriod)
	if err != nil {
		errLog.Fatal(err)
	}
	//defers
	defer rdb.Client.Close()
	defer wg.Wait()
	// initialize api handler
	cli := &http.Client{} // TODO check the bad practice
	api := sms_ru.NewAPIHandler(
		*smsRUAPIID,
		cli,
		infLog,
		errLog,
		sms_ru.WithTest(),
		sms_ru.JSONFormat(),
	)
	// initialize application
	myApp := sms_service.Application{
		ErrorLog:     errLog,
		InfoLog:      infLog,
		APIsms:       api,
		SendInterval: *sendInterval,
		DB:           rdb,
		JobsWG:       &wg,
	}
	myApp.RunJobs()

	// initialize server
	server := &http.Server{
		Addr:     serverPort,
		ErrorLog: errLog,
		Handler:  myApp.Router(),
	}
	go server.ListenAndServe()
	infLog.Printf("Run server - port: %s", serverPort)

	s := <-sigChan
	server.Close()
	infLog.Printf("Finish service, signal: %s", s)
}
