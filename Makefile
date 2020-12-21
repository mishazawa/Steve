include .env
export

build:
	go build -o bin/$(BASE) .

run:
	go run .

commit:
	gofmt -w .
	git add .
	git commit --allow-empty
