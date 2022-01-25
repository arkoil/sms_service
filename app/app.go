package app

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type Application struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	SmsListDB *redis.Client
	CTX       *context.Context
}
