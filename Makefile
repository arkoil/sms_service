NAME_WEB=web_service

build:
	go build -o bin/${NAME_WEB} cmd/web_service/main.go
run_web:
	./bin/${NAME_WEB}
clean:
	go clean
	rm bin/${NAME_WEB}
