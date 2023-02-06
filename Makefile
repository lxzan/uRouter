PWD=$(shell pwd)

tidy:
	cd ${PWD}/contrib/adapter/gws && go mod tidy
	cd ${PWD}/contrib/adapter/http && go mod tidy
	go mod tidy

test:
	cd ${PWD}/contrib/adapter/gws && go test --count=1 ./...
	cd ${PWD}/contrib/adapter/gorilla && go test --count=1 ./...
	cd ${PWD}/contrib/adapter/http && go test --count=1 ./...
	go test --count=1 ./...

cover:
	go test -coverprofile=./bin/cover.out --cover ./...
