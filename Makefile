.PHONY: all binary build default test 

default: all

all: binary

get:
		go get ./...

test: get 
		go test ./... -v 


binary: get
	GOOS=windows GOARCH=amd64 go build -o ipmes2yammer.exe main.go 


