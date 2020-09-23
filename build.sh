#!/usr/bin/env bash

GOOS=linux go build -o ./bin/kotaplan main.go
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o ./bin/kotaplan-darwin main.go
GOOS=windows GOARCH=386 go build -o ./bin/kotaplan.exe
shasum -a 256 bin/kotaplan bin/kotaplan-darwin bin/kotaplan.exe > ./bin/sha256
