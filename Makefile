include .env
export

commit:
	gofmt -w .
	git add .
	git commit --allow-empty

run-bot:
	go run ./bot/

run-runner:
	go run ./runner/ -remote http://$(HOST):$(PORT) -local ./data

run-runner-upload:
	go run ./runner/ -remote http://$(HOST):$(PORT) -local ./data -upload

run-waiter:
	go run -race ./waiter/ -port $(PORT) -datafolder ./temp

build-bot:
	go build -o bin/$(BASE)-bot ./bot/

build-bot-arm:
	GOARM=7	GOARCH=arm GOOS=linux go build -o bin/$(BASE)-bot-arm ./bot/

build-runner:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BASE)-runner-win ./runner/
	GOOS=linux   GOARCH=amd64 go build -o bin/$(BASE)-runner-linux ./runner/

build-waiter-arm:
	GOARM=7	GOARCH=arm GOOS=linux go build -o bin/$(BASE)-waiter-arm ./waiter/

build:
	make build-bot
	make build-bot-arm
	make build-runner
	make build-waiter-arm
