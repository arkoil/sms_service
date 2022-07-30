# sms service

### Run command example
```
$ .bin/bin sms_service -rdb=8 -sms-ru-aoi-id="API_ID"
```
#### Option
`-interval` interval between sending SMS in seconds (default 60)

`-port string` set service port (default "7575")


`-prefbe string` sets the Redis key for sent messages (default "sms:besend:")

`-prefto` sets the Redis key to send messages (default "sms:tosend:")

`-rdb set` redis db (default 11)

`-rpass` set redis password

`-rurl` set redis url (default "localhost:6379")

`-sp` set storage period of sent SMS history in hours (default 720)

`-sms-ru-aoi-id` set sms.ru apiID

### Endpoints

#### /sms/send

##### Params:

`request_id` - your message id

`number` - phone number

`message` - your message text

#### /sms/check

##### Params:

`request_id` - your message id


### Reference

#### Used API`s
https://sms.ru - sms send api
