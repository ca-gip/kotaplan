#!/usr/bin/env bash

GOOS=linux go build -o ./bin/rqc main.go
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o ./bin/rqc-darwin main.go
GOOS=windows GOARCH=386 go build -o ./bin/rqc.exe
md5sum ./bin/rqc ./bin/rqc.exe ./bin/rqc-darwin > ./bin/md5
