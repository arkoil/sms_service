package store

import (
	"context"
	"encoding/json"
	"github.com/arkoil/sms_service/internal/sms"
	"github.com/go-redis/redis/v8"
	"time"
)

type DB struct {
	Client          *redis.Client
	keyPrefixToSend string
	keyPrefixBeSend string
	storePeriod     time.Duration
}

type SMS interface {
	Phone() string
	Message() string
	ID() string
}

var Ctx = context.TODO()

func NewDB(addr string, pass string, dbNum int, keyPrefixToSend string, keyPrefixBeSend string, storePeriod int) (*DB, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbNum,
	})
	if err := cli.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &DB{
		Client:          cli,
		keyPrefixToSend: keyPrefixToSend,
		keyPrefixBeSend: keyPrefixBeSend,
		storePeriod:     time.Duration(storePeriod) * time.Hour,
	}, nil

}

func (d DB) GetSMSList() ([]SMS, error) {
	keys, err := d.Client.Keys(Ctx, d.keyPrefixToSend+"*").Result()
	if err != nil {
		return nil, err
	}
	smsList := make([]SMS, 0, len(keys))
	for _, key := range keys {
		jsonItem, _ := d.Client.Get(Ctx, key).Result()
		item := sms.Item{}
		err := json.Unmarshal([]byte(jsonItem), &item)
		if err != nil {
			return smsList, err
		}
		smsList = append(smsList, item)
		d.Client.Del(Ctx, key)
	}
	return smsList, nil
}
func (d DB) SetSMSSent(id string, status string) error {
	key := d.keyPrefixBeSend + id
	err := d.Client.Set(Ctx, key, status, d.storePeriod).Err()
	if err != nil {
		return err
	}
	return nil
}
