package main

import (
	"context"
	"flag"
	"github.com/arkoil/sms_service/internal/web_service"
	"github.com/arkoil/sms_service/internal/web_service/tasks"
	"github.com/go-redis/redis/v8"
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
	flag.Parse()

	//Variables
	ctx := context.Background()
	serverPort := strings.Join([]string{":", *port}, "")
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	var wg sync.WaitGroup
	//Clean close
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	//Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       11, // use default DB
	})

	//defers
	defer rdb.Close()
	defer wg.Wait()

	myApp := web_service.Application{
		errLog,
		infLog,
		rdb,
		&ctx,
	}
	server := &http.Server{
		Addr:     serverPort,
		ErrorLog: errLog,
		Handler:  myApp.Router(),
	}

	infLog.Printf("Запуск сервера на %s", serverPort)
	go tasks.SNSSender(rdb, &ctx, errLog, infLog, &wg)
	go server.ListenAndServe()
	//errLog.Fat al(err)
	s := <-sigChan

	server.Close()
	infLog.Printf("Завершение по сигналу: %s", s)
}
