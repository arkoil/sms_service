NAME_SERVICE=sms_service

build:
	go build -o bin/${NAME_SERVICE} cmd/${NAME_SERVICE}/main.go
clean:
	go clean
	rm bin/${NAME_SERVICE}
