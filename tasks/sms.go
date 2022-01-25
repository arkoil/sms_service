package tasks

import (
	"context"
	"encoding/json"
	"github.com/arkoil/sms_service/sms"
	"github.com/arkoil/sms_service/sms_ru"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const sleepTime = 1 * time.Minute

//const callbackURL = "http://lk.kuskovo.biz"
const callbackURL = "http://localhost:5757"

func SNSSender(db *redis.Client, ctx *context.Context, errLog *log.Logger, infLog *log.Logger, wg *sync.WaitGroup) {
	defer SNSSender(db, ctx, errLog, infLog, wg)
	time.Sleep(sleepTime)
	infLog.Print("Начинаем отправку смс")
	wg.Add(1)
	defer wg.Done()
	keys, err := db.Keys(*ctx, "sms:sending:*").Result()
	if err != nil {
		errLog.Println(err)
	} else {
		cli := &http.Client{}
		api := sms_ru.NewAPIHandler(
			"FC9FD058-2026-9AF1-34DA-2471603BE87F",
			cli,
			true,
		)
		smsList := make([]sms_ru.SMS, 0)
		for _, key := range keys {
			infLog.Println(key)
			jsonItem, _ := db.Get(*ctx, key).Result()
			item := sms.Item{}
			err := json.Unmarshal([]byte(jsonItem), &item)
			if err != nil {
				errLog.Print(err)
				continue
			}
			smsList = append(smsList, item)
			db.Del(*ctx, key)

		}
		res, err := api.SendSMSList(&smsList, true)
		if err != nil {
			errLog.Print(err)
		} else {
			parseRes := sms_ru.ApiResponse{}
			err = json.Unmarshal([]byte(res), &parseRes)
			if err != nil {
				errLog.Print(err)
			} else {
				infLog.Print(parseRes)
				for number, data := range parseRes.SMS {
					for _, sms := range smsList {
						if sms.Phone() == number {
							params := url.Values{}
							params.Set("request_id", sms.ID())
							params.Set("sms_api_id", data.SMSID)
							params.Set("status_code", string(data.StatusCode))
							url := callbackURL + "/api/request/check/sms"
							req, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
							if err != nil {
								errLog.Print(err)
							} else {
								infLog.Println(req.Body)
							}
						}
					}
				}

			}
		}
	}
	infLog.Print("Завершаем отправку смс")
}
